async function getApiKey(): Promise<string | undefined> {
  return await fetch('/refresh', { method: 'GET' })
    .catch(() => undefined)
    .then((resp) => resp?.json().then((body) => body.token as string));
}

export const theme = () => ({
  authToken: getApiKey(),
});
