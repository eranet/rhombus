# Example Gazebo RhoMBus integration

## Dependencies

### CNats

- `git clone https://github.com/nats-io/cnats.git`
- `mkdir build && cd build`
- `cmake .. -DNATS_BUILD_STREAMING=OFF && make`
- `sudo make install`
- `sudo ldconfig /usr/local/lib`

`sudo apt install libgazebo9-dev`

## Building

- `mkdir build && cd build`
- `cmake .. && make`

## Testing

- Launch nats server
- `./launch_world.sh gripper.world`
- `go run sub/joint_state.go`
