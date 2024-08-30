const os = require('os');
const fs = require('fs');
const path = require('path');

const platform = os.platform();
let bin = '';
let targetPath = path.join(__dirname, 'bin', 'makeCsr');

if (platform === 'win32') {
  bin = 'makeCsr-windows-amd64.exe';
  targetPath += '.exe';  // 在 Windows 上目标文件名加上 .exe 扩展名
} else if (platform === 'darwin') {
  bin = 'makeCsr-darwin-amd64';
} else if (platform === 'linux') {
  bin = 'makeCsr-linux-amd64';
} else {
  console.error(`Unsupported platform: ${platform}`);
  process.exit(1);
}

try {
  const sourcePath = path.join(__dirname, 'bin', bin);

  // 检查文件是否存在
  if (!fs.existsSync(sourcePath)) {
    throw new Error(`Binary file not found: ${sourcePath}`);
  }

  // 复制文件并设置权限
  fs.copyFileSync(sourcePath, targetPath);

  // Windows 不需要设置执行权限，跳过 chmod
  if (platform !== 'win32') {
    fs.chmodSync(targetPath, 0o755);
  }

  console.log('Binary copied successfully');
} catch (err) {
  console.error(`Failed to copy binary: ${err.message}`);
  process.exit(1);
}
