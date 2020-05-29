
# coding:utf-8

# pip install fastapi
# pip install uvicorn  #需要一个ASGI服务器来进行生产


from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()


class Item(BaseModel):
    name: str
    price: float
    is_offer: bool = None


@app.get("/ping")
def read_root():
    return {"ping": "pong"}


# url /items/1000?q=value
@app.get("/items/{item_id}")
def read_item(item_id: int, q: str = None):
    return {"item_id": item_id, "q": q}


@app.put('/items/{item_id}')
async def update_item(item_id: int, item: Item):
    return {'item_id': item_id, 'item_name': item.name}
