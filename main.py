import os
import sys
import subprocess
import json
import logging
import hashlib
import math
from time import sleep

from froModuleDrivers.nioManager import NioManager as IoManager
from froModuleDrivers.mCarDriver import MCarDriver as CarDriver

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(filename)s - %(threadName)s: [%(levelname)s] - %("
    "message)s",
)
logger = logging.getLogger(__name__)

REPLAN = False

if __name__ == "__main__":
    root = os.path.dirname(__file__)
    if not os.path.exists(os.path.join(root, "bin", "main")):
        logger.info("Building...")
        subprocess.run(["/bin/sh", os.path.join(root, "build.sh")]).check_returncode()
        logger.info("Build complete!")
    with open(os.path.join(root, "config", "plan.json")) as f:
        md5hash = hashlib.md5(f.read().encode("utf-8")).hexdigest()
    result_filename = os.path.join(root, "result", f"{md5hash}.json")
    if not os.path.exists(result_filename) or REPLAN:
        logger.info("Planning...")
        subprocess.run(
            [
                os.path.join(root, "bin", "main"),
                os.path.join(root, "config", "plan.json"),
                result_filename,
            ]
        ).check_returncode()
        logger.info("Route planning complete!")
    logger.info("Loading config...")
    with open(result_filename) as f:
        result_data: dict = json.load(f)
    with open(os.path.join(root, "config", "car.json")) as f:
        cat_data: dict = json.load(f)
    logger.info("Config loaded!")
    car = CarDriver(cat_data["host"], int(cat_data.get("port", 4001)))
    io_manager = IoManager(car)
    io_manager.run()
    for _ in range(3):
        if car.refreshStatus():
            break
    else:
        logger.info("Failed to start car!")
        sys.exit(1)
    logger.info("地图：\n" + result_data["map"])
    logger.info("路线：\n" + result_data["string"])
    logger.info("Starting...")
    try: 
        for path in result_data["path"]:
            for card in path["left_turning_points"]:
                car.setTurnLeftCard(card)
            for card in path["right_turning_points"]:
                car.setTurnRightCard(card)
            car.parkPoint(path["dest"])
            car.startCar()
            card_id = car.getCardNum()
            while not card_id or card_id != path["dest"]:
                sleep(0.2)
            sleep(math.max(path["goods_num"] - 1, 0))
            if path["goods_num"] != 0:
                car.turnLeftOrigin()
    except KeyboardInterrupt:
        logger.error("Interrupted!")
        sys.exit(1)
    except ConnectionResetError:
        logger.error("Connection reset!")
        sys.exit(1)
    else:
        logger.info("Done!")
