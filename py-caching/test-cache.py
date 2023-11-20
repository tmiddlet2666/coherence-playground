from coherence import NamedCache, Session
import asyncio
import requests
import time

session: Session
cache: NamedCache[int, str]


async def init_coherence():
    global session
    session = await Session.create()

    global cache
    cache = await session.get_cache("test")


def coherence_cache(func):
    async def wrapper(*args):
        if await cache.contains_key(args):
            return await cache.get(args)
        else:
            result = func(*args)
            await cache.put(args, result)
            return result

    return wrapper


@coherence_cache
def get_html_data(url):
    response = requests.get(url)
    return response.text


async def run_test():
    global cache
    await init_coherence()

    start_time = time.time()
    result = await get_html_data('https://www.oracle.com/')
    print('Time taken - first call: ', time.time() - start_time)


    start_time = time.time()
    result = await get_html_data('https://www.oracle.com/')
    print('Time taken - second call: ', time.time() - start_time)

    # await cache.put(1, "hello")
    #
    # # get the value for a key in the cache
    # r = await cache.get(1)
    #
    # # print the value got for a key in the cache
    # print("The value of key 1 is " + r)

asyncio.run(run_test())