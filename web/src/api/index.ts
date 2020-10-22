export default class API {
  static test(val: string) {
    // eslint-disable-next-line
    console.log(val, process.env.NODE_ENV, process.env.BACKEND_URL);
  }

  query = async (input: any):Promise<any> => {
    const res = await fetch(process.env.BACKEND_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
      body: JSON.stringify(input),
    });
    if (res.status !== 200) {
      return { error: 'error' };
    }
    return res.json();
  }
}
