#include "joint_state_publisher.h"

//TODO: REFACTORME!

using namespace gazebo;

JointStatePublisher::JointStatePublisher() {}
JointStatePublisher::~JointStatePublisher()
{
	if (this->nats_conn) {
		natsConnection_Destroy(this->nats_conn);
	}
}

static void OnMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure)
{
    auto x = json::parse(natsMsg_GetData(msg));
    physics::Joint* joint = (physics::Joint*) closure;
    joint->SetPosition(0,x["Value"], true);

    // Don't forget to destroy the message!
    natsMsg_Destroy(msg);
}

void JointStatePublisher::Load(physics::ModelPtr _parent, sdf::ElementPtr _sdf)
{
	this->parent_ = _parent;
	this->world_ = _parent->GetWorld();

	// publish topic
	std::string robotNS = parent_->GetName();
	if (_sdf->HasElement("robotNamespace")) {
		robotNS = _sdf->GetElement("robotNamespace")->Get<std::string>();
		if (robotNS.empty()) robotNS = parent_->GetName();
	}

	this->publish_topic = robotNS + "/joint_states";

	std::cout << "Publish topic: " << this->publish_topic << "\n";
	
	// update rate
	this->update_rate_ = 1000.0;
	if (_sdf->HasElement("updateRate")) {
		this->update_rate_ = _sdf->GetElement("updateRate")->Get<double>();
	}
	if (this->update_rate_ > 0.0) {
		this->update_period_ = 1.0 / this->update_rate_;
	} else {
		this->update_period_ = 0.0;
	}
	this->last_update_time_ = this->world_->SimTime();

	// Nats connect
	natsStatus s;
	s = natsConnection_ConnectTo(&this->nats_conn, NATS_DEFAULT_URL);

	// bind on update
	this->updateConnection = event::Events::ConnectWorldUpdateBegin(
			std::bind(&JointStatePublisher::OnUpdate, this));

    std::string name = "";
    std::vector<physics::JointPtr> joints = this->parent_->GetJoints();
    for (int i = 0; i < joints.size(); i++) {
        name = joints[i]-> GetName();
        natsSubscription    *sub = NULL;
        //natsConnection_Subscribe(&sub, this->nats_conn, (robotNS+"/"+ name+"/command").c_str(),OnMsg , &(joints[i]));
        natsConnection_Subscribe(&sub, this->nats_conn, (robotNS+"/"+ name+"/command").c_str(),OnMsg , joints[i].get());
    }
}

void JointStatePublisher::OnUpdate()
{
	common::Time current_time = this->world_->SimTime();
	if (current_time < this->last_update_time_)
	{
		this->last_update_time_ = current_time;
	}

	double seconds_since_last_update = (current_time - last_update_time_).Double();

	if (seconds_since_last_update > this->update_period_) {
		publishJointStates();
		this->last_update_time_+= common::Time(update_period_);
	}
}

void JointStatePublisher::publishJointStates()
{
	std::vector<physics::JointPtr> joints = this->parent_->GetJoints();

	std::vector<std::string> names ;
	std::vector<double> poss ;
	std::vector<double> vels ;
    for (int i = 0; i < joints.size(); i++) {
        names.push_back(joints[i]->GetName());
        poss.push_back(joints[i]->Position(0));
        vels.push_back(joints[i]->GetVelocity(0));
    }

    json obj;
    obj["name"] = names;
    obj["position"] = poss;
    obj["velocity"] = vels;

    std::string msg = obj.dump();

	natsStatus s = natsConnection_PublishString(this->nats_conn,
			this->publish_topic.c_str(), msg.c_str());

	//std::cout << "Publish msg: Done." << "\n";
}
