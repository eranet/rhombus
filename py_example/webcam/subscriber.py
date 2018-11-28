import asyncio
from datetime import datetime
from nats.aio.client import Client as NATS


async def run(loop):
    nc = NATS()

    await nc.connect(loop=loop)

    async def message_handler(msg):
        with open('received.jpg', 'wb') as f:
            f.write(msg.data)

    await nc.subscribe('webcam', cb=message_handler, is_async=True)    


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    loop.run_forever()
    loop.close()