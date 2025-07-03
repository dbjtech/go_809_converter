// API基础URL
//获取当前url完整路径
let url = window.location.href;
//从 /static处截断，前面部分就是根路径
const staticIndex = url.indexOf('/static');
let API_BASE = '';
// 如果 staticIndex 等于 -1 则返回当前文件的目录,即最后一个 / 所在位置
if (staticIndex === -1) {
    API_BASE = url.substring(0, url.lastIndexOf('/'));
} else {
    API_BASE = url.substring(0, staticIndex);
}
// DOM元素
const toolStatusDiv = document.getElementById('toolStatus');
const fileSelectBtn = document.getElementById('fileSelectBtn');
const fileNameLabel = document.getElementById('file-name-label');
const fileSelectText = document.getElementById('fileSelectText');
const fileDropdown = document.getElementById('fileDropdown');
const selectAllBtn = document.getElementById('selectAllBtn');
const clearAllBtn = document.getElementById('clearAllBtn');
const fileInfo = document.getElementById('fileInfo');
const searchPattern = document.getElementById('searchPattern');
const searchBtn = document.getElementById('searchBtn');
const clearBtn = document.getElementById('clearBtn');
const cancelBtn = document.getElementById('cancelBtn');
const searchResults = document.getElementById('searchResults');
const loading = document.getElementById('loading');

// 全局变量
let selectedFiles = new Set();
let allFiles = [];
let currentSearchController = null; // 用于取消搜索的AbortController
let currentStreamReader = null; // 用于取消流式搜索的Reader
let streamAbort = false;

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function () {
    checkToolStatus().then(r => console.log("检查工具状态成功"));
    loadLogFiles().then(r => console.log("加载日志列表成功"));
    setupEventListeners();
});

// 设置事件监听器
function setupEventListeners() {
    searchBtn.addEventListener('click', performSearch);
    clearBtn.addEventListener('click', () => { searchPattern.value = '' });
    cancelBtn.addEventListener('click', cancelSearch);
    selectAllBtn.addEventListener('click', selectAllFiles);
    clearAllBtn.addEventListener('click', clearAllFiles);
    fileNameLabel.addEventListener('dblclick', function () {
        loadLogFiles().then(r => console.log("刷新日志列表成功"));
        clearAllFiles();
    });

    // 添加流式搜索按钮事件
    const streamSearchBtn = document.getElementById('streamSearchBtn');
    if (streamSearchBtn) {
        streamSearchBtn.addEventListener('click', performStreamSearch);
    }

    // 阻止多选下拉框自动收起
    fileDropdown.addEventListener('click', function (e) {
        // 阻止点击复选框时收起下拉框
        if (e.target.type === 'checkbox' || e.target.closest('.form-check')) {
            e.stopPropagation();
        }
    });

    // 点击下拉框外部时收起
    document.addEventListener('click', function (e) {
        if (!fileSelectBtn.contains(e.target) && !fileDropdown.contains(e.target)) {
            // 使用Bootstrap的方法收起下拉框
            const dropdown = bootstrap.Dropdown.getInstance(fileSelectBtn);
            if (dropdown) {
                dropdown.hide();
            }
        }
    });

    // 回车键搜索
    searchPattern.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            performSearch();
        }
    });
}

// 检查工具状态
async function checkToolStatus() {
    try {
        const response = await fetch(`${API_BASE}/tools/check`);
        const data = await response.json();

        displayToolStatus(data);
    } catch (error) {
        console.error('检查工具状态失败:', error);
        toolStatusDiv.innerHTML = `
            <div class="alert alert-danger">
                <i class="fas fa-exclamation-triangle"></i>
                无法检测系统工具状态: ${error.message}
            </div>
        `;
    }
}

// 显示工具状态
function displayToolStatus(toolData) {
    const tools = [
        { name: 'cat', label: 'Cat', available: toolData.has_cat },
        { name: 'gzip', label: 'Gzip', available: toolData.has_gzip },
        { name: 'grep', label: 'Grep', available: toolData.has_grep }
    ];

    let html = '';

    tools.forEach(tool => {
        const statusClass = tool.available ? 'tool-available' : 'tool-unavailable';
        const icon = tool.available ? 'fas fa-check' : 'fas fa-times';

        html += `
            <span class="tool-status ${statusClass}">
                <i class="${icon}"></i> ${tool.label}
            </span>
        `;
    });

    // 添加处理策略指示器
    const hasAllTools = toolData.has_cat && toolData.has_gzip && toolData.has_grep;
    const strategyIcon = hasAllTools ? 'fas fa-rocket text-success' : 'fas fa-cog text-warning';
    const strategyTitle = hasAllTools ? '高性能模式' : '兼容模式';

    html += `
        <span class="tool-status" style="background-color: transparent; border: none; color: #6c757d;" title="${hasAllTools ? '使用系统工具处理' : '使用Go内置方法'}">
            <i class="${strategyIcon}"></i> ${strategyTitle}
        </span>
    `;

    toolStatusDiv.innerHTML = html;
}

// 加载日志文件列表
async function loadLogFiles() {
    try {
        const response = await fetch(`${API_BASE}/logs/list`);
        const data = await response.json();

        populateFileSelect(data.files);
    } catch (error) {
        console.error('加载文件列表失败:', error);
        fileDropdown.innerHTML = '<li class="file-item"><div class="alert alert-danger m-2">加载文件列表失败</div></li>';
    }
}

// 填充文件选择下拉框
function populateFileSelect(files) {
    allFiles = files;

    // 清空现有的文件项（保留头部和控制按钮）
    const existingItems = fileDropdown.querySelectorAll('.file-item');
    existingItems.forEach(item => item.remove());

    files.forEach(file => {
        const li = document.createElement('li');
        li.className = 'file-item';

        const compressionIcon = file.compressed ?
            '<i class="fas fa-file-archive text-info"></i>' :
            '<i class="fas fa-file-alt text-secondary"></i>';

        li.innerHTML = `
            <div class="form-check px-3 py-1">
                <input class="form-check-input" type="checkbox" value="${file.id}" id="file_${file.id.replace(/[^a-zA-Z0-9]/g, '_')}" data-file-info='${JSON.stringify(file)}'>
                <label class="form-check-label w-100" for="file_${file.id.replace(/[^a-zA-Z0-9]/g, '_')}" style="font-size: 12px;">
                    <div class="d-flex align-items-center">
                        ${compressionIcon}
                        <span class="ms-2 flex-grow-1">${file.name}</span>
                        <small class="text-muted">${formatFileSize(file.size)}</small>
                    </div>
                </label>
            </div>
        `;

        // 添加复选框变化事件
        const checkbox = li.querySelector('input[type="checkbox"]');
        checkbox.addEventListener('change', handleFileSelection);

        fileDropdown.appendChild(li);
    });
}

// 处理文件选择
function handleFileSelection(event) {
    const checkbox = event.target;
    const fileId = checkbox.value;

    if (checkbox.checked) {
        selectedFiles.add(fileId);
    } else {
        selectedFiles.delete(fileId);
    }

    updateFileSelectionDisplay();
}

// 全选文件
function selectAllFiles() {
    const checkboxes = fileDropdown.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.checked = true;
        selectedFiles.add(checkbox.value);
    });
    updateFileSelectionDisplay();
}

// 清空选择
function clearAllFiles() {
    const checkboxes = fileDropdown.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.checked = false;
    });
    selectedFiles.clear();
    updateFileSelectionDisplay();
}

// 更新文件选择显示
function updateFileSelectionDisplay() {
    const count = selectedFiles.size;

    if (count === 0) {
        fileSelectText.textContent = '选择文件...';
        fileInfo.innerHTML = '';
    } else if (count === 1) {
        const selectedFileId = Array.from(selectedFiles)[0];
        const fileData = allFiles.find(f => f.id === selectedFileId);
        if (fileData) {
            fileSelectText.textContent = fileData.name;
            const compressionIcon = fileData.compressed ?
                '<i class="fas fa-file-archive text-info"></i>' :
                '<i class="fas fa-file-alt text-secondary"></i>';
            fileInfo.innerHTML = `
                <div class="d-flex align-items-center" style="font-size: 11px;">
                    ${compressionIcon}
                    <span class="ms-1 text-muted">${formatFileSize(fileData.size)} | ${fileData.modified}</span>
                </div>
            `;
        }
    } else {
        fileSelectText.textContent = `已选择 ${count} 个文件`;
        const totalSize = Array.from(selectedFiles).reduce((sum, fileId) => {
            const fileData = allFiles.find(f => f.id === fileId);
            return sum + (fileData ? fileData.size : 0);
        }, 0);
        fileInfo.innerHTML = `
            <div class="d-flex align-items-center" style="font-size: 11px;">
                <i class="fas fa-files text-primary"></i>
                <span class="ms-1 text-muted">总计: ${formatFileSize(totalSize)}</span>
            </div>
        `;
    }
}

// 执行搜索
function performSearch() {
    searchLogs();
}

// 执行流式搜索
function performStreamSearch() {
    streamSearchLogs();
}

// 取消搜索
function cancelSearch() {
    // 取消普通搜索
    if (currentSearchController) {
        currentSearchController.abort();
        currentSearchController = null;
    }

    // 取消流式搜索
    if (currentStreamReader) {
        currentStreamReader.cancel();
        currentStreamReader = null;
        streamAbort = true;
    }

    // 清理待渲染队列（如果存在）
    if (typeof pendingResults !== 'undefined') {
        pendingResults.length = 0;
        isRenderScheduled = false;
    }

    // 重置UI状态
    resetSearchUI();
}

// 重置搜索UI状态
function resetSearchUI() {
    loading.style.display = 'none';
    searchBtn.disabled = false;
    const streamSearchBtn = document.getElementById('streamSearchBtn');
    if (streamSearchBtn) {
        streamSearchBtn.disabled = false;
    }
    cancelBtn.style.disabled = 'disabled';
}

// 搜索日志
function searchLogs() {
    const pattern = searchPattern.value.trim();

    if (!pattern) {
        alert('请输入搜索关键词');
        return;
    }

    if (selectedFiles.size === 0) {
        alert('请选择要搜索的日志文件');
        return;
    }

    // 创建AbortController用于取消搜索
    currentSearchController = new AbortController();

    // 显示加载状态
    loading.style.display = 'block';
    searchResults.innerHTML = '';
    searchBtn.disabled = true;
    const streamSearchBtn = document.getElementById('streamSearchBtn');
    if (streamSearchBtn) {
        streamSearchBtn.disabled = true;
    }
    cancelBtn.style.disabled = '';

    // 如果只选择了一个文件，使用原有的单文件搜索API
    if (selectedFiles.size === 1) {
        const selectedFileId = Array.from(selectedFiles)[0];
        const fileData = allFiles.find(f => f.id === selectedFileId);
        if (!fileData) {
            alert('文件信息未找到');
            loading.style.display = 'none';
            searchBtn.disabled = false;
            return;
        }

        const requestData = {
            file_id: fileData.id,  // 使用文件ID而不是路径
            pattern: pattern
        };

        fetch(`${API_BASE}/logs/search`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData),
            signal: currentSearchController.signal
        })
            .then(response => response.json())
            .then(data => {
                currentSearchController = null;
                resetSearchUI();

                if (data.error) {
                    searchResults.innerHTML = `<div class="alert alert-danger">${escapeHtml(data.error)}</div>`;
                } else {
                    displaySearchResults(data.results, pattern, [selectedFileId]);
                }
            })
            .catch(error => {
                currentSearchController = null;
                resetSearchUI();

                if (error.name === 'AbortError') {
                    searchResults.innerHTML = `<div class="alert alert-warning">搜索已取消</div>`;
                } else {
                    searchResults.innerHTML = `<div class="alert alert-danger">搜索失败: ${escapeHtml(error.message)}</div>`;
                }
            });
    } else {
        // 多文件搜索：并行搜索所有选中的文件
        searchMultipleFiles(Array.from(selectedFiles), pattern);
    }
}

// 流式搜索日志
function streamSearchLogs() {
    const pattern = searchPattern.value.trim();

    if (!pattern) {
        alert('请输入搜索关键词');
        return;
    }

    if (selectedFiles.size === 0) {
        alert('请选择要搜索的日志文件');
        return;
    }

    // if (selectedFiles.size > 1) {
    //     alert('流式搜索目前只支持单个文件');
    //     return;
    // }

    // 显示加载状态
    loading.style.display = 'block';
    searchResults.innerHTML = '';
    searchBtn.disabled = true;
    const streamSearchBtn = document.getElementById('streamSearchBtn');
    if (streamSearchBtn) {
        streamSearchBtn.disabled = true;
    }

    const selectedFileId = Array.from(selectedFiles);
    const selectedFileData = allFiles.filter(f => selectedFileId.includes(f.id));
    if (!selectedFileData) {
        alert('文件信息未找到');
        loading.style.display = 'none';
        searchBtn.disabled = false;
        if (streamSearchBtn) {
            streamSearchBtn.disabled = false;
        }
        return;
    }
    // 初始化结果容器和渲染节流变量
    let resultCount = 0;
    let pendingResults = [];
    let isRenderScheduled = false;
    const RENDER_BATCH_SIZE = 100; // 每次最多渲染100个结果
    const RENDER_INTERVAL = 15; // 最小渲染间隔15ms
    for (let fileIndex = 0; fileIndex < selectedFileData.length; fileIndex++) {
        const fileData = selectedFileData[fileIndex];
        if (streamAbort) {
            break;
        }
        searchResults.innerHTML = `
        <div class="search-summary mb-3">
            <div class="d-flex justify-content-between align-items-center">
                <span class="text-muted">
                    <i class="fas fa-search"></i>
                    实时搜索中... <span id="resultCounter">找到 <strong>0</strong> 个匹配项</span>
                </span>
                <button class="btn btn-sm btn-outline-secondary" onclick="clearResults()">
                    <i class="fas fa-times"></i> 清除结果
                </button>
            </div>
        </div>
        <div class="search-results-container" id="streamResults">
            <div class="alert alert-info">
                <i class="fas fa-spinner fa-spin"></i> 正在搜索文件: ${escapeHtml(fileData.name)}
            </div>
        </div>
    `;

        const streamResults = document.getElementById('streamResults');
        const resultCounter = document.getElementById('resultCounter');

        // 搜索完成后的清理函数
        function finishSearchCleanup() {
            resetSearchUI();

            // 移除加载提示
            const loadingAlert = streamResults.querySelector('.alert-info');
            if (loadingAlert) {
                loadingAlert.remove();
            }

            if (resultCount === 0) {
                streamResults.innerHTML = '<div class="alert alert-info">未找到匹配的结果</div>';
            }
        }

        // 批量渲染函数
        function renderPendingResults() {
            // 检查是否已取消搜索, 如果缓存大于1000，取消则停止渲染
            if (streamAbort && pendingResults.length > 1000) {
                pendingResults.length = 0;
                isRenderScheduled = false;
                streamAbort = false;
                return;
            }

            if (pendingResults.length === 0) {
                isRenderScheduled = false;
                return;
            }

            // 移除加载提示（如果还存在）
            const loadingAlert = streamResults.querySelector('.alert-info');
            if (loadingAlert) {
                loadingAlert.remove();
            }

            // 批量处理结果，每次最多处理RENDER_BATCH_SIZE个
            const batchToRender = pendingResults.splice(0, RENDER_BATCH_SIZE);
            // 如果页面有内容，并且待处理的结果数量大于2.1万，则忽略中间的数据，只处理最后面的20000个结果
            if (streamResults.children.length > 2000 && pendingResults.length > 21000) {
                setTimeout(() => {
                    renderPendingResults();
                }, 0);
                return;
            }
            const fragment = document.createDocumentFragment();

            batchToRender.forEach(resultData => {
                const highlightedLine = highlightMatches(resultData.content, resultData.matched);
                const resultDiv = document.createElement('div');
                resultDiv.className = 'search-result-item mb-2';
                resultDiv.innerHTML = `
                <div class="result-header">
                    <span class="line-number">
                        <i class="fas fa-hashtag"></i> 第 ${resultData.line_number} 行
                    </span>
                    <span class="match-info text-muted">
                        匹配 ${resultData.matchIndex}
                    </span>
                </div>
                <div class="result-content">
                    <code class="result-line">${highlightedLine}</code>
                </div>
            `;
                fragment.appendChild(resultDiv);
            });

            // 检查子元素数量，如果超过21000就删除到20000左右
            if (streamResults.children.length > 21000) {
                const elementsToRemove = streamResults.children.length - 20000;
                for (let i = 0; i < elementsToRemove; i++) {
                    if (streamResults.firstElementChild) {
                        streamResults.removeChild(streamResults.firstElementChild);
                    }
                }
            }

            // 一次性添加所有元素到DOM
            streamResults.appendChild(fragment);

            // 更新计数器
            let tips = resultCount > 20000 ? "。匹配大于 2w 条。请稍等结果展示，这里只展示最后的约20000 条结果。" : "";
            resultCounter.innerHTML = `找到 <strong>${resultCount}</strong> 个匹配项${tips}`;

            // 优化滚动行为：减少滚动频率，避免性能问题
            if (batchToRender.length > 0) {
                // 只在结果数量较少时或每50个结果时才滚动
                if (resultCount <= 50 || resultCount % 50 === 0) {
                    const lastResult = streamResults.lastElementChild;
                    if (lastResult && lastResult.classList.contains('search-result-item')) {
                        // 使用instant滚动而不是smooth，避免动画累积
                        lastResult.scrollIntoView({ behavior: 'instant', block: 'nearest' });
                    }
                }
            }

            // 如果还有待渲染的结果，继续调度下一次渲染
            if (pendingResults.length > 0) {
                setTimeout(() => {
                    renderPendingResults();
                }, RENDER_INTERVAL);
            } else {
                isRenderScheduled = false;
            }
        }

        // 创建EventSource连接
        const requestData = {
            file_id: fileData.id,
            pattern: pattern
        };

        function parseMultilineString(str) {
            const lines = str.split('&nl ');
            const result = [];
            let currentItem = null;

            lines.forEach(line => {
                // 匹配标签行（格式：标签名:值）
                const labelMatch = line.match(/^([^:\s]+):\s*(.*)$/);

                if (labelMatch) {
                    // 是标签行，创建新对象
                    const [, label, value] = labelMatch;
                    currentItem = [label, value];
                    result.push(currentItem);
                } else if (currentItem) {
                    // 不是标签行，将当前行追加到上一个值中
                    currentItem[1] += '\n' + line;
                }
            });

            return result;
        }
        // 使用fetch发送POST请求到SSE端点
        fetch(`${API_BASE}/logs/search-stream`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const reader = response.body.getReader();
                currentStreamReader = reader;
                const decoder = new TextDecoder();
                let buffer = ''; // 添加缓冲区来处理跨chunk的数据

                function readStream() {
                    return reader.read().then(({ done, value }) => {
                        if (done) {
                            // 处理缓冲区中剩余的数据
                            if (buffer.trim()) {
                                processBufferedData(buffer);
                            }

                            // 搜索完成
                            currentStreamReader = null;
                            resetSearchUI();

                            // 移除加载提示
                            const loadingAlert = streamResults.querySelector('.alert-info');
                            if (loadingAlert) {
                                loadingAlert.remove();
                            }

                            if (resultCount === 0) {
                                streamResults.innerHTML = '<div class="alert alert-info">未找到匹配的结果</div>';
                            }
                            return;
                        }

                        const chunk = decoder.decode(value);
                        buffer += chunk; // 将新数据添加到缓冲区

                        // 按双换行符分割SSE消息块
                        const messages = buffer.split('\n\n');

                        // 保留最后一个可能不完整的消息在缓冲区中
                        buffer = messages.pop() || '';

                        // 处理完整的消息
                        for (const message of messages) {
                            if (message.trim()) {
                                processSSEMessage(message);
                            }
                        }

                        return readStream();
                    });
                }

                // 处理单个SSE消息
                function processSSEMessage(message) {
                    const lines = message.split('\n');
                    let currentEvent = null;
                    let currentData = null;

                    for (const line of lines) {
                        if (line.startsWith('event: ')) {
                            currentEvent = line.substring(7).trim();
                        } else if (line.startsWith('data: ')) {
                            currentData = line.substring(6);
                        }
                    }

                    // 处理完整的事件数据
                    if (currentEvent && currentData !== null) {
                        processEventData(currentEvent, currentData);
                    }
                }

                // 处理缓冲区中剩余的数据
                function processBufferedData(data) {
                    if (data.trim()) {
                        processSSEMessage(data);
                    }
                }

                // 处理事件数据
                function processEventData(eventType, data) {
                    if (eventType === 'result' && data) {
                        try {
                            // 解析字符串
                            const parsedLines = parseMultilineString(data);
                            if (parsedLines.length < 3) {
                                console.error('解析搜索结果格式错误:', data);
                                return;
                            }
                            const line_number = parsedLines[0][1];
                            const matched = parsedLines[1][1];
                            const content = parsedLines[2][1];

                            if (line_number && content) {
                                // 这是一个搜索结果
                                resultCount++;

                                // 添加到待渲染队列而不是立即渲染
                                pendingResults.push({
                                    line_number: line_number,
                                    content: content,
                                    matched: matched,
                                    matchIndex: resultCount
                                });

                                // 如果没有调度渲染，则启动渲染
                                if (!isRenderScheduled) {
                                    isRenderScheduled = true;
                                    console.log("启动渲染");
                                    // 使用setTimeout而不是requestAnimationFrame，确保在后台标签页也能正常工作
                                    setTimeout(renderPendingResults, RENDER_INTERVAL);
                                }
                            }
                        } catch (e) {
                            console.error('解析搜索结果失败:', e, data);
                        }
                    } else if (eventType === 'finished') {
                        // 搜索完成，停止处理
                        currentStreamReader = null;

                        // 确保所有待渲染的结果都被处理完毕
                        if (pendingResults.length > 0) {
                            // 异步处理剩余结果，避免阻塞主线程
                            const finishRendering = () => {
                                if (pendingResults.length > 0) {
                                    renderPendingResults();
                                    // 继续异步处理剩余结果
                                    setTimeout(finishRendering, 0);
                                } else {
                                    // 所有结果渲染完成后的清理工作
                                    finishSearchCleanup();
                                }
                            };
                            finishRendering();
                        } else {
                            // 没有待渲染结果，直接清理
                            finishSearchCleanup();
                        }

                        // 结束流读取
                        if (currentStreamReader) {
                            currentStreamReader.cancel();
                        }
                        return;
                    } else if (eventType === 'error') {
                        try {
                            const errorData = JSON.parse(data);
                            throw new Error(errorData.error || '搜索过程中发生错误');
                        } catch (e) {
                            throw new Error('搜索过程中发生未知错误');
                        }
                    }
                }

                return readStream();
            })
            .catch(error => {
                currentStreamReader = null;
                resetSearchUI();

                if (error.name === 'AbortError') {
                    searchResults.innerHTML = `<div class="alert alert-warning">搜索已取消</div>`;
                } else {
                    searchResults.innerHTML = `<div class="alert alert-danger">流式搜索失败: ${escapeHtml(error.message)}</div>`;
                }
            });

    }
}

// 搜索多个文件
function searchMultipleFiles(files, pattern) {
    const searchPromises = files.map(fileId => {
        const fileData = allFiles.find(f => f.id === fileId);
        if (!fileData) {
            return Promise.resolve({ file: fileId, error: '文件信息未找到' });
        }

        const requestData = {
            file_id: fileData.id,  // 使用文件ID而不是路径
            pattern: pattern
        };

        return fetch(`${API_BASE}/logs/search`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData),
            signal: currentSearchController.signal
        })
            .then(response => response.json())
            .then(data => ({ file: fileId, data }))
            .catch(error => ({ file: fileId, error: error.message }));
    });

    Promise.all(searchPromises)
        .then(results => {
            currentSearchController = null;
            resetSearchUI();

            // 合并所有结果
            const allResults = [];
            const errors = [];

            results.forEach(result => {
                if (result.error) {
                    errors.push(`${result.file}: ${result.error}`);
                } else if (result.data.error) {
                    errors.push(`${result.file}: ${result.data.error}`);
                } else if (result.data.results) {
                    // 为每个结果添加文件信息
                    result.data.results.forEach(match => {
                        allResults.push({
                            ...match,
                            sourceFile: result.file
                        });
                    });
                }
            });

            if (errors.length > 0) {
                searchResults.innerHTML = `<div class="alert alert-warning">部分文件搜索失败:<br>${errors.map(e => escapeHtml(e)).join('<br>')}</div>`;
            }

            if (allResults.length > 0) {
                displaySearchResults(allResults, pattern, files);
            } else if (errors.length === 0) {
                searchResults.innerHTML = '<div class="alert alert-info">未找到匹配的结果</div>';
            }
        })
        .catch(error => {
            currentSearchController = null;
            resetSearchUI();

            if (error.name === 'AbortError') {
                searchResults.innerHTML = `<div class="alert alert-warning">搜索已取消</div>`;
            } else {
                searchResults.innerHTML = `<div class="alert alert-danger">搜索失败: ${escapeHtml(error.message)}</div>`;
            }
        });
}

// 显示搜索结果
function displaySearchResults(results, pattern, searchedFiles = []) {
    if (!results || results.length === 0) {
        searchResults.innerHTML = `
            <div class="alert alert-info">
                <i class="fas fa-info-circle"></i>
                未找到匹配的结果
            </div>
        `;
        return;
    }

    const isMultiFile = searchedFiles.length > 1;
    const totalMatches = results.length;

    let html = `
        <div class="search-summary mb-3">
            <div class="d-flex justify-content-between align-items-center">
                <span class="text-muted">
                    <i class="fas fa-search"></i>
                    找到 <strong>${totalMatches}</strong> 个匹配项
                    ${isMultiFile ? `，搜索了 <strong>${searchedFiles.length}</strong> 个文件` : ''}
                </span>
                <button class="btn btn-sm btn-outline-secondary" onclick="clearResults()">
                    <i class="fas fa-times"></i> 清除结果
                </button>
            </div>
        </div>
        <div class="search-results-container">
    `;

    // 如果是多文件搜索，按文件分组显示
    if (isMultiFile) {
        const resultsByFile = {};
        results.forEach(result => {
            const file = result.sourceFile || 'unknown';
            if (!resultsByFile[file]) {
                resultsByFile[file] = [];
            }
            resultsByFile[file].push(result);
        });

        Object.keys(resultsByFile).forEach(file => {
            const fileResults = resultsByFile[file];
            const fileName = file.split('/').pop() || file.split('\\').pop() || file;

            html += `
                <div class="file-group mb-4">
                    <div class="file-header bg-light p-2 rounded" style="cursor: pointer; user-select: none;" onclick="toggleFileResults(this)">
                        <h6 class="mb-0">
                            <i class="fas fa-chevron-down toggle-icon me-2"></i>
                            <i class="fas fa-file-alt text-primary"></i>
                            ${escapeHtml(fileName)}
                            <span class="badge bg-primary ms-2">${fileResults.length} 个匹配</span>
                        </h6>
                        <small class="text-muted">${escapeHtml(file)}</small>
                    </div>
                    <div class="file-results mt-2">
            `;

            fileResults.forEach((result, index) => {
                const highlightedLine = highlightMatches(result.content, pattern);
                html += `
                    <div class="search-result-item mb-2 ms-3">
                        <div class="result-header">
                            <span class="line-number">
                                <i class="fas fa-hashtag"></i> 第 ${result.line_number} 行
                            </span>
                            <span class="match-info text-muted">
                                匹配 ${index + 1}
                            </span>
                        </div>
                        <div class="result-content">
                            <code class="result-line">${highlightedLine}</code>
                        </div>
                    </div>
                `;
            });

            html += `
                    </div>
                </div>
            `;
        });
    } else {
        // 单文件搜索，直接显示结果
        results.forEach((result, index) => {
            const highlightedLine = highlightMatches(result.content, pattern);

            html += `
                <div class="search-result-item mb-2">
                    <div class="result-header">
                        <span class="line-number">
                            <i class="fas fa-hashtag"></i> 第 ${result.line_number} 行
                        </span>
                        <span class="match-info text-muted">
                            匹配 ${index + 1}
                        </span>
                    </div>
                    <div class="result-content">
                        <code class="result-line">${highlightedLine}</code>
                    </div>
                </div>
            `;
        });
    }

    html += '</div>';
    searchResults.innerHTML = html;
}

// 高亮匹配的文本
function highlightMatches(content, pattern) {
    // 转义HTML特殊字符
    const escapedContent = escapeHtml(content);
    const escapedPattern = pattern;

    // 使用正则表达式高亮匹配项
    const regex = new RegExp(`(${escapedPattern})`, 'gi');
    return escapedContent.replace(regex, '<span class="matched-text">$1</span>');
}

// 清空搜索结果
function clearResults() {
    searchResults.innerHTML = '';
    // searchPattern.value = '';
    // 清空文件选择
    // selectedFiles.clear();
    // updateFileSelectionDisplay();
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';

    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 转义HTML特殊字符
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 转义正则表达式特殊字符
function escapeRegExp(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

// 切换文件结果的显示/隐藏
function toggleFileResults(headerElement) {
    const fileGroup = headerElement.parentElement;
    const fileResults = fileGroup.querySelector('.file-results');
    const toggleIcon = headerElement.querySelector('.toggle-icon');

    if (fileResults.style.display === 'none') {
        // 展开
        fileResults.style.display = 'block';
        toggleIcon.className = 'fas fa-chevron-down toggle-icon me-2';
        headerElement.style.backgroundColor = '';
    } else {
        // 收起
        fileResults.style.display = 'none';
        toggleIcon.className = 'fas fa-chevron-right toggle-icon me-2';
        headerElement.style.backgroundColor = '#f8f9fa';
    }
}