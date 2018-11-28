import asyncio
from datetime import datetime

from cv2 import VideoCapture, imencode, imwrite
from nats.aio.client import Client as NATS


async def run(loop):
    nc = NATS()
    await nc.connect(io_loop=loop)   

    cam = VideoCapture(0)
    while True:
        ok, img = cam.read()
        if ok:
            ok, buf = imencode('.jpg', img)
            if ok:    
                print(str(datetime.now()))
                await nc.publish('webcam', buf)
                await asyncio.sleep(1)

    await nc.close()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))    
    loop.close()
