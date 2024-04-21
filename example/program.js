result.addMetadata("foo", "bar");
result.setMetadata("foo", "bar2");
result.addMetadata("foo", "bar3");
result.setConfig(
  "topic",
  `http://localhost:8989/${message.topic()}?foo=bar&bar=baz`
);
result.setConfig("method", "POST");

const foo = cache.get("foo") ?? 0;
const bar = cache.get("bar") ?? "";

const time = new Date().toISOString();

result.setData({
  time,
  data: message.topic(),
  foo,
  bar,
});

cache.set("foo", foo + 1);
cache.set("bar", `bar: ${time}`);
