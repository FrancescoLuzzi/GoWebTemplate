type callback = (args: unknown) => void;
interface subscriber {
  token: number;
  func: callback;
}
export interface PubSub {
  sub: (topic: string, callback: callback) => number;
  pub: (topic: string, value: unknown) => boolean;
  unsub: (id: number) => boolean;
}

const pubsub = {
  topics: {} as Record<string, subscriber[]>,
  subUid: -1,
  sub(topic: string, func: callback): number {
    if (!this.topics[topic]) {
      this.topics[topic] = [];
    }
    const token = ++this.subUid;
    this.topics[topic].push({
      token: token,
      func: func,
    });
    return token;
  },
  pub(topic: string, args: unknown): boolean {
    if (!this.topics[topic]) {
      return false;
    }
    const subscribers = this.topics[topic];
    setTimeout(function () {
      let len = subscribers ? subscribers.length : 0;
      while (len--) {
        subscribers[len].func(args);
      }
    }, 0);
    return true;
  },
  unsub(token: number): boolean {
    for (const topic in this.topics) {
      if (this.topics[topic]) {
        for (let i = 0, j = this.topics[topic].length; i < j; i++) {
          if (this.topics[topic][i].token === token) {
            this.topics[topic].splice(i, 1);
            return true;
          }
        }
      }
    }
    return false;
  },
};

export default pubsub as PubSub;
