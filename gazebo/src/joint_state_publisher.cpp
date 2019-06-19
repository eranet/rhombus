#include "joint_state_publisher.h"
#include <boost/algorithm/string.hpp>


using namespace gazebo;

JointStatePublisher::JointStatePublisher() {}
JointStatePublisher::~JointStatePublisher()
{
	if (this->nats_conn) {
		natsConnection_Destroy(this->nats_conn);
	}
}

void JointStatePublisher::Load(physics::ModelPtr _parent, sdf::ElementPtr _sdf)
{
	this->parent_ = _parent;
	this->world_ = _parent->GetWorld();

	// publish topic
	this->publish_topic = parent_->GetName();
	if (_sdf->HasElement("robotNamespace")) {
		this->publish_topic = _sdf->GetElement("robotNamespace")->Get<std::string>();
		if (this->publish_topic.empty()) this->publish_topic = parent_->GetName();
	}
	this->publish_topic += "/joint_states";

	std::cout << "Publish topic: " << this->publish_topic << "\n";
	


	// update rate
	this->update_rate_ = 100.0;
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
	std::vector<physics::JointPtr> joints_ = this->parent_->GetJoints();
	std::string name = "[";
	std::string position = "[";
	std::string velocity = "[";
	for (int i = 0; i < joints_.size(); i++) {
		name += "\"" + joints_[i]->GetName() + "\", ";
		position += std::to_string(joints_[i]->Position(0)) + ", ";
		velocity += std::to_string(joints_[i]->Position(0)) + ", ";
	}
	this->parent_->GetJoint("palm_left_finger")->SetVelocity(0, 1.0);
	name.resize(name.size()-2);
	position.resize(position.size()-2);
	velocity.resize(velocity.size()-2);
	name += "]";
	position += "]";
	velocity += "]";
	std::string msg = "{\"name\": " + name + ", \"position\": " + position +
		", \"velocity\": " + velocity + "}";

	// Nats publish
	std::cout << "Publish msg: " << msg << "\n";

	natsStatus s = natsConnection_PublishString(this->nats_conn,
			this->publish_topic.c_str(), msg.c_str());

	std::cout << "Publish msg: Done." << "\n";
}
