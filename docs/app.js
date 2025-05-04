const modeButtons = document.querySelectorAll(".mode-btn");
const urlPanel = document.querySelector(".panel-url");
const repoPanel = document.querySelector(".panel-repo");
const urlInput = document.getElementById("urlInput");
const outputInput = document.getElementById("outputInput");
const repoInput = document.getElementById("repoInput");
const pathsInput = document.getElementById("pathsInput");
const branchInput = document.getElementById("branchInput");
const outputInput2 = document.getElementById("outputInput2");
const commandOutput = document.getElementById("commandOutput");
const message = document.getElementById("message");
const copyBtn = document.getElementById("copyBtn");
const chips = document.querySelectorAll(".chip");

let mode = "url";

function setMode(next) {
  mode = next;
  modeButtons.forEach((btn) => btn.classList.toggle("is-active", btn.dataset.mode === next));
  urlPanel.classList.toggle("is-active", next === "url");
  repoPanel.classList.toggle("is-active", next === "repo");
  updateCommand();
}

function sanitize(value) {
  return value.trim();
}

function parseGitHubUrl(raw) {
  try {
    const url = new URL(raw);
    if (url.hostname !== "github.com") {
      return { ok: false, error: "Only github.com URLs are supported." };
    }

    const parts = url.pathname.split("/").filter(Boolean);
    if (parts.length < 5) {
      return { ok: false, error: "URL must include /tree/<branch>/path or /blob/<branch>/path." };
    }

    const [owner, repoRaw, mode, branch, ...rest] = parts;
    const repo = repoRaw.replace(/\.git$/, "");

    if (!owner || !repo || !branch || rest.length === 0) {
      return { ok: false, error: "Missing owner, repo, branch, or path." };
    }

    if (mode !== "tree" && mode !== "blob") {
      return { ok: false, error: "URL must be a tree or blob link." };
    }

    let subPath = rest.join("/");
    if (mode === "blob") {
      subPath = rest.slice(0, -1).join("/");
      if (!subPath) {
        return { ok: false, error: "Blob URL must point inside a directory." };
      }
    }

    return {
      ok: true,
      repoUrl: `https://github.com/${owner}/${repo}`,
      branch,
      subPath,
    };
  } catch (err) {
    return { ok: false, error: "Invalid URL." };
  }
}

function buildCommand() {
  if (mode === "url") {
    const url = sanitize(urlInput.value);
    const out = sanitize(outputInput.value);
    if (!url) {
      return { cmd: "gitsub clone …", msg: "Paste a GitHub directory or file URL." };
    }

    const parsed = parseGitHubUrl(url);
    if (!parsed.ok) {
      return { cmd: "gitsub clone …", msg: parsed.error };
    }

    const outputArg = out ? ` -o ${out}` : "";
    return {
      cmd: `gitsub clone ${url}${outputArg}`,
      msg: `Detected: ${parsed.repoUrl} · ${parsed.branch} · ${parsed.subPath}`,
    };
  }

  const repo = sanitize(repoInput.value);
  const paths = sanitize(pathsInput.value);
  const branch = sanitize(branchInput.value);
  const out = sanitize(outputInput2.value);

  if (!repo || !paths) {
    return { cmd: "gitsub clone …", msg: "Provide a repo URL and one or more directories." };
  }

  const branchArg = branch ? ` -b ${branch}` : "";
  const outputArg = out ? ` -o ${out}` : "";
  return { cmd: `gitsub clone ${repo} ${paths}${branchArg}${outputArg}`, msg: "" };
}

function updateCommand() {
  const { cmd, msg } = buildCommand();
  commandOutput.textContent = cmd;
  message.textContent = msg;
}

function copyCommand() {
  const text = commandOutput.textContent;
  if (!text || text.includes("…")) return;

  navigator.clipboard.writeText(text).then(
    () => {
      copyBtn.textContent = "Copied!";
      setTimeout(() => {
        copyBtn.textContent = "Copy";
      }, 1200);
    },
    () => {
      copyBtn.textContent = "Copy failed";
      setTimeout(() => {
        copyBtn.textContent = "Copy";
      }, 1200);
    }
  );
}

modeButtons.forEach((btn) => btn.addEventListener("click", () => setMode(btn.dataset.mode)));
urlInput.addEventListener("input", updateCommand);
outputInput.addEventListener("input", updateCommand);
repoInput.addEventListener("input", updateCommand);
pathsInput.addEventListener("input", updateCommand);
branchInput.addEventListener("input", updateCommand);
outputInput2.addEventListener("input", updateCommand);
copyBtn.addEventListener("click", copyCommand);

chips.forEach((chip) =>
  chip.addEventListener("click", () => {
    setMode("url");
    urlInput.value = chip.dataset.example;
    updateCommand();
  })
);

updateCommand();
