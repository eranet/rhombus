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

- `mkdir -p build && cd build`
- `cmake .. && make`

## Testing

- Launch nats server
- `./launch_world.sh`
- `go run sub/joint_state.go`
- `go run pub/publish_cmd.go`

## Conversion
gz sdf -p my_urdf.urdf > my_sdf.sdf