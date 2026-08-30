function timeAgo(epochSeconds) {
  const diff = Date.now() / 1000 - epochSeconds;
  const units = [
    ["year", 31536000],
    ["month", 2592000],
    ["day", 86400],
    ["hour", 3600],
    ["minute", 60],
  ];

  for (const [name, secs] of units) {
    const value = Math.floor(diff / secs);
    if (value >= 1) return `${value} ${name}${value > 1 ? "s" : ""} ago`;
  }

  return "just now";
}

function renderTimeAgo() {
  document.querySelectorAll("time.timeago[data-ts]").forEach((el) => {
    const ts = parseFloat(el.dataset.ts);
    if (!isNaN(ts) && ts > 0) {
      el.textContent = timeAgo(ts);
    }
  });
}

document.addEventListener("DOMContentLoaded", renderTimeAgo);
document.body.addEventListener("htmx:afterSwap", renderTimeAgo);
setInterval(renderTimeAgo, 120_000);
