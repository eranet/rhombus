#include <gazebo/gazebo.hh>
#include <gazebo/physics/physics.hh>
#include <gazebo/common/common.hh>
#include <nats/nats.h>

#include <json.hpp>
using json = nlohmann::json;

namespace gazebo {
class JointStatePublisher : public ModelPlugin {
public:
	JointStatePublisher();
	~JointStatePublisher();
	void Load(physics::ModelPtr _parent, sdf::ElementPtr _sdf);
	void OnUpdate();
	void publishJointStates();

	physics::ModelPtr parent_;
private:
	event::ConnectionPtr updateConnection;
	physics::WorldPtr world_;


	// Update Rate
	double update_rate_;
	double update_period_;
	common::Time last_update_time_;

	// Nats
	natsConnection *nats_conn;
	std::string publish_topic;
};

GZ_REGISTER_MODEL_PLUGIN (JointStatePublisher)
}
