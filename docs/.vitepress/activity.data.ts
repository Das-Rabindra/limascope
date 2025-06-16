declare const data: { stars: number; pulls: number };
export { data };

export default {
  async load() {
    const urls = [
      "https://api.github.com/repos/Das-Rabindra/limascope",
      "https://hub.docker.com/v2/namespaces/Das-Rabindra/limascope",
    ];

    const responses = await Promise.all(urls.map((url) => fetch(url).then((res) => res.json())));

    const data = {
      stars: responses[0].stargazers_count,
      pulls: responses[1].pull_count,
    };

    return data;
  },
};
