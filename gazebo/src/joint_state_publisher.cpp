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

typedef struct
{
  physics::JointControllerPtr jointController;
  std::string name;
}ParamStruct;


static void OnMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure)
{
    auto x = json::parse(natsMsg_GetData(msg));
    //std::cout << "1: \n" ;
    ParamStruct* p = (ParamStruct*) closure;
    //std::cout << "2: \n" ;
    //physics::JointControllerPtr joint = p->joint;
    //std::cout << "3: \n" ;
    //std::cout << p->name << "4: \n" ;
    p->jointController->SetPositionTarget(p->name,x["Value"]);
    //free(closure);
    // Don't forget to destroy the message!
    natsMsg_Destroy(msg);
    //std::cout << "5: \n" ;
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
	s = natsConnection_ConnectTo(&this->nats_conn, NATS_DEFAULT_URL); //TODO: default url

	// bind on update
	this->updateConnection = event::Events::ConnectWorldUpdateBegin(
			std::bind(&JointStatePublisher::OnUpdate, this));

    physics::JointControllerPtr jc = this->parent_->GetJointController();
    std::map<std::string, physics::JointPtr> mjc = jc->GetJoints();
    for( auto const& [key, val] : mjc )
{
    std::cout << key         // string (key)
              << ':'
              << val->GetName()        // string's value
              << std::endl ;
}

    std::string name = "";
    std::vector<physics::JointPtr> joints = this->parent_->GetJoints();
    for (int i = 0; i < joints.size(); i++) {
        name = joints[i]-> GetName();
        natsSubscription    *sub = NULL;
        std::cout << "Name: " << joints[i]->GetName() << "\n";
        std::cout << "Scoped Name: " << joints[i]->GetScopedName() << "\n";
        //natsConnection_Subscribe(&sub, this->nats_conn, (robotNS+"/"+ name+"/command").c_str(),OnMsg , &(joints[i]));
        ParamStruct* p = new ParamStruct();
        //std::cout << "Crea: ";
        p->jointController = jc;
        //std::cout << "set1: ";
        p->name = joints[i]->GetScopedName();
        //std::cout << "set2: ";
        //std::tuple<physics::JointPtr, std::string> param = std::make_tuple(joints[i].get(), joints[i]->GetScopedName());
        natsConnection_Subscribe(&sub, this->nats_conn, (robotNS+"/"+ name+"/command").c_str(),OnMsg , p);

        gazebo::common::PID pid = gazebo::common::PID(1000,1,0.1,100,0.0,-1.0,0.0);
        jc->SetPositionPID(joints[i]-> GetScopedName(), pid);
    }
    jc->Update();

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
