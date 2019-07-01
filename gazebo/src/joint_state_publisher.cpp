#include "joint_state_publisher.h"
//#include <boost/algorithm/string.hpp>


//TODO: REFACTORME! FIX THIS FUCKING C

using namespace gazebo;

JointStatePublisher::JointStatePublisher() {}
JointStatePublisher::~JointStatePublisher()
{
	if (this->nats_conn) {
		natsConnection_Destroy(this->nats_conn);
	}
}

//typedef std::pair<JointStatePublisher*, std::string> p ;

static void OnMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure)
{

    // Prints the message, using the message getters:
//    printf("Received msg: %s - %.*s\n",
//        natsMsg_GetSubject(msg),
//        natsMsg_GetDataLength(msg),
//        natsMsg_GetData(msg));


    //Cut middle from: simple_gripper/right_finger_tip/command
    std::string subj = natsMsg_GetSubject(msg) ;
    subj = subj.substr(0, subj.rfind("/"));
    std::string name = subj.substr(subj.rfind("/")+1);
    std::cout <<  name ;

    auto x = json::parse(natsMsg_GetData(msg));
    std::cout << " Command: "<< x["Value"] << std::endl;

    JointStatePublisher* jsp =  (JointStatePublisher*) closure; //SHIT

    physics::JointPtr joint = jsp->parent_->GetJoint(name);
    joint->SetPosition(0,x["Value"], true);

    std::cout << "Command Done" << std::endl;

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
	this->update_rate_ = 10.0;
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
        natsConnection_Subscribe(&sub, this->nats_conn, (robotNS+"/"+ name+"/command").c_str(),OnMsg , this);
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
	// Hardcoded JSON serialization
	std::vector<physics::JointPtr> joints = this->parent_->GetJoints();

	std::vector<std::string> names ;
	std::vector<double> poss ;
	std::vector<double> vels ;
    for (int i = 0; i < joints.size(); i++) {
        names.push_back(joints[i]->GetName());
        poss.push_back(joints[i]->Position(0));
        vels.push_back(joints[i]->GetVelocity(0));
    }


//	std::string name = "[";
//	std::string position = "[";
//	std::string velocity = "[";
//	for (int i = 0; i < joints_.size(); i++) {
//		name += "\"" + joints_[i]->GetName() + "\", ";
//		position += std::to_string(joints_[i]->Position(0)) + ", ";
//		velocity += std::to_string(joints_[i]->Position(0)) + ", ";
//	}
	//this->parent_->GetJoint("palm_left_finger")->SetVelocity(0, 1.0);


	//sim_joints_[j]->SetPosition(0, joint_position_command_[j], true);

//	name.resize(name.size()-2);
//	position.resize(position.size()-2);
//	velocity.resize(velocity.size()-2);
//	name += "]";
//	position += "]";
//	velocity += "]";
//	std::string msg = "{\"name\": " + name + ", \"position\": " + position +
//		", \"velocity\": " + velocity + "}";

    json obj;
    obj["name"] = names;
    obj["position"] = poss;
    obj["velocity"] = vels;

    std::string msg = obj.dump();

	// Nats publish
	//std::cout << "Publish msg: " << msg << "\n";

	natsStatus s = natsConnection_PublishString(this->nats_conn,
			this->publish_topic.c_str(), msg.c_str());

	//std::cout << "Publish msg: Done." << "\n";
}
