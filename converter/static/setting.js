/*
 * @Author: SimingLiu siming.liu@dbjtech.com
 * @Date: 2025-01-30 12:00:00
 * @LastEditors: SimingLiu siming.liu@dbjtech.com
 * @LastEditTime: 2025-02-19 10:50:00
 * @FilePath: \go_809_converter\converter\static\setting.js
 * @Description: 配置管理页面JavaScript逻辑
 * 
 */
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
// 全局变量
let currentConfig = {};
let originalConfig = {};
let modifiedKeys = new Set();

// 工具函数：将 dot-keys 对象还原为嵌套结构
function unflattenToNested(flat) {
    const nested = {};
    Object.keys(flat).forEach(k => {
        const parts = k.split('.');
        let cur = nested;
        for (let i = 0; i < parts.length; i++) {
            const part = parts[i];
            if (i === parts.length - 1) {
                cur[part] = flat[k];
            } else {
                if (!cur[part] || typeof cur[part] !== 'object') {
                    cur[part] = {};
                }
                cur = cur[part];
            }
        }
    });
    return nested;
}

// 工具函数：将嵌套对象扁平化为 dot-keys（缺失补回）
function flattenObject(obj, prefix = '') {
    const result = {};
    if (obj && typeof obj === 'object' && !Array.isArray(obj)) {
        Object.keys(obj).forEach((key) => {
            const val = obj[key];
            const newKey = prefix ? `${prefix}.${key}` : key;
            if (val && typeof val === 'object' && !Array.isArray(val)) {
                Object.assign(result, flattenObject(val, newKey));
            } else {
                result[newKey] = val;
            }
        });
    }
    return result;
}

// 页面加载完成后初始化
$(document).ready(function() {
    initializePage();
    loadCurrentConfig();
    loadConfigTree();
});

// 值比较辅助函数：根据原始类型和key进行规范化比较
function valuesEqualByKey(key, originalValue, rawInputValue) {
    const isCrypto = typeof key === 'string' && (key.endsWith('.cryptoPacket') || key === 'cryptoPacket');
    if (isCrypto) {
        const toArray = (v) => Array.isArray(v)
            ? v.map(s => String(s).trim()).filter(s => s.length > 0)
            : typeof v === 'string'
                ? v.split(',').map(s => s.trim()).filter(s => s.length > 0)
                : [];
        const curr = toArray(rawInputValue);
        const orig = toArray(originalValue);
        if (curr.length !== orig.length) return false;
        for (let i = 0; i < curr.length; i++) {
            if (curr[i] !== orig[i]) return false;
        }
        return true;
    }
    const origType = typeof originalValue;
    if (origType === 'number') {
        const num = parseFloat(rawInputValue);
        // 如果原始是数字，严格比较数值；非数字视为不相等
        return Number.isFinite(num) && num === originalValue;
    }
    if (origType === 'boolean') {
        const bool = rawInputValue === 'true';
        return bool === originalValue;
    }
    // 默认按字符串比较（去除null/undefined影响）
    const origStr = originalValue == null ? '' : String(originalValue);
    const curStr = rawInputValue == null ? '' : String(rawInputValue);
    return origStr === curStr;
}

// 页面加载完成后初始化
$(document).ready(function() {
    initializePage();
    loadCurrentConfig();
    loadConfigTree();
});

// 初始化页面
function initializePage() {
    // 标签页切换事件
    $('.tab').click(function() {
        const tabId = $(this).data('tab');
        switchTab(tabId);
    });

    // 监听配置项变化（按类型精确比较，避免改回原值仍被标注）
    $(document).on('input change', '.config-input', function() {
        const key = $(this).data('key');
        const originalValue = originalConfig[key];
        const rawInput = $(this).val();
        
        const equal = valuesEqualByKey(key, originalValue, rawInput);
        
        if (!equal) {
            modifiedKeys.add(key);
            $(this).addClass('modified');
            $(this).closest('.config-item').find('label').addClass('modified');
        } else {
            modifiedKeys.delete(key);
            $(this).removeClass('modified');
            $(this).closest('.config-item').find('label').removeClass('modified');
        }
        
        updateConfigStatus();
    });
}

// 获取当前环境名称（优先使用 env 字段）
function getCurrentEnvironment() {
    try {
        const env = currentConfig && typeof currentConfig['env'] === 'string' ? currentConfig['env'].trim() : '';
        if (env) return env;
        const keys = Object.keys(currentConfig || {});
        const candidates = ['develop', 'online', 'staging', 'production'];
        for (const cand of candidates) {
            if (keys.some(k => k.startsWith(cand + '.'))) return cand;
        }
        const firstTop = keys.map(k => k.split('.')[0]).find(seg => seg && seg !== 'env');
        return firstTop || 'develop';
    } catch (e) {
        return 'develop';
    }
}

// 渲染顶部特殊行：env、jtwTcp、normalTcp 同行展示；其中 jtwTcp、normalTcp 为只读复选框，均不可删除
function renderSpecialTopRow() {
    const $container = $('#config-sections');
    const envVal = (currentConfig['env'] ?? '').toString();
    const jtwVal = !!currentConfig['jtwTcp'];
    const normalVal = !!currentConfig['normalTcp'];

    const $section = $('<div class="config-section">');
    $section.append('<h3>启动命令参数</h3>');

    const $row = $('<div class="inline-row">');
    $row.css({ display: 'flex', alignItems: 'center', gap: '100px', flexWrap: 'wrap' });

    // env：可编辑文本输入，参与保存（用于环境迁移）
    const $envItem = $('<div class="config-item">');
    $envItem.css({ display: 'flex', alignItems: 'center', gap: '8px' });
    $envItem.append('<label>env:</label>');
    $envItem.append(`<input type="text" class="config-input" data-key="env" value="${envVal}" readonly>`);

    // jtwTcp：只读复选框（不参与保存，不可删除）
    const jtwKey = 'jtwTcp';
    const $jtwItem = $('<div class="config-item">');
    $jtwItem.css({ display: 'flex', alignItems: 'center', gap: '0px' });
    $jtwItem.append(`<label>${DISPLAY_LABELS[jtwKey] || jtwKey}:</label>`);
    $jtwItem.append(`<input type="checkbox" ${jtwVal ? 'checked' : ''} onclick="return false">`);

    // normalTcp：只读复选框（不参与保存，不可删除）
    const normalKey = 'normalTcp';
    const $normalItem = $('<div class="config-item">');
    $normalItem.css({ display: 'flex', alignItems: 'center', gap: '8px' });
    $normalItem.append(`<label>${DISPLAY_LABELS[normalKey] || normalKey}:</label>`);
    $normalItem.append(`<input type="checkbox" ${normalVal ? 'checked' : ''} onclick="return false">`);

    $row.append($envItem, $jtwItem, $normalItem);
    $section.append($row);
    $container.append($section);
}

// 渲染配置区域
function renderConfigSections() {
    const $container = $('#config-sections');
    $container.empty();

    // 先渲染顶部特殊行
    renderSpecialTopRow();
    
    // 按配置组分组（跳过顶级 env/jtwTcp/normalTcp）
    const groups = {};
    Object.keys(currentConfig).forEach(key => {
        if (key === 'env' || key === 'jtwTcp' || key === 'normalTcp') {
            return; // 顶级特殊项不进入普通分组
        }
        const parts = key.split('.');
        const groupName = parts.length > 2 ? parts.slice(0, -1).join('.') : parts[0];
        
        if (!groups[groupName]) {
            groups[groupName] = {};
        }
        groups[groupName][key] = currentConfig[key];
    });
    
    const envName = getCurrentEnvironment();
    // 渲染每个配置组
    Object.keys(groups).forEach(groupName => {
        const $section = $('<div class="config-section">');
        // 仅对转换连接组（env.converter.project）提供删除按钮
        const isSubProject = groupName.startsWith(envName + '.converter.') && groupName.split('.').length === 3;
        const headerHtml = `<h3>${groupName} ${isSubProject ? `<button class="btn btn-danger btn-sm" onclick="deleteConfigGroup('${groupName}')" style="margin-left: 10px;">删除转换连接</button>` : ''}</h3>`;
        $section.append(headerHtml);
        
        if (isSubProject) {
            // 转换连接配置：enabled 和 name 优先，其他按字母顺序
            const sortedKeys = sortSubProjectKeys(Object.keys(groups[groupName]));
            sortedKeys.forEach(key => {
                const value = groups[groupName][key];
                const $item = createConfigItem(key, value);
                $section.append($item);
            });
        } else {
            // 其他配置：按原有逻辑（不显示删除按钮在子项上）
            Object.keys(groups[groupName]).forEach(key => {
                const value = groups[groupName][key];
                const $item = createConfigItem(key, value);
                $section.append($item);
            });
        }
        
        $container.append($section);
    });
}

// 对转换连接配置键进行排序
function sortSubProjectKeys(keys) {
    const priorityFields = ['enable', 'name'];
    const priorityKeys = [];
    const otherKeys = [];
    
    keys.forEach(key => {
        const fieldName = key.split('.').pop();
        if (priorityFields.includes(fieldName)) {
            priorityKeys.push(key);
        } else {
            otherKeys.push(key);
        }
    });
    
    // enabled 和 name 按优先级排序
    priorityKeys.sort((a, b) => {
        const fieldA = a.split('.').pop();
        const fieldB = b.split('.').pop();
        return priorityFields.indexOf(fieldA) - priorityFields.indexOf(fieldB);
    });
    
    // 其他字段按字母顺序排序
    otherKeys.sort((a, b) => {
        const fieldA = a.split('.').pop();
        const fieldB = b.split('.').pop();
        return fieldA.localeCompare(fieldB);
    });
    
    return [...priorityKeys, ...otherKeys];
}

// 创建配置项（普通项，移除子项级删除按钮）
// 中文展示名称映射
const DISPLAY_LABELS = {
  name: '上级名称',
  enable: '是否开启本连接',
  enabled: '是否开启本连接',
  cryptoPacket: '需加密的推送报文',
  encryptKey: '加密密钥',
  extendVersion: '是否使用DBJ扩展后的809协议',
  govServerIP: '上级平台ip',
  govServerPort: '上级平台端口',
  jtw809ConverterDownLinkIp: '暴露给交委下行链接的ip',
  jtw809ConverterDownLinkPort: '暴露给交委下行连接的端口',
  jtw809ConverterIp: '交委上级平台ip',
  jtw809ConverterPort: '交委上级平台端口',
  localServerIP: '暴露给上级平台下行链接的ip',
  localServerPort: '暴露给上级平台下行连接的端口',
  openCrypto: '开启加密',
  platformId: '上级平台连接id',
  platformPassword: '上级平台连接密码',
  platformUserId: '上级平台分配的用户ID',
  protocolVersion: '协议版本',
  thirdpartPort: '暴露给第三方推送连接的端口',
  useLocationInterval: '1分钟内最多推送一个位置点',
  database: '数据库名',
  host: '数据库地址',
  password: '数据库密码',
  pool_idle_conns: '数据库空闲连接数',
  pool_size: '数据库连接池',
  port: '数据库端口',
  showSQL: '是否打印sql日志',
  user: '数据库连接用户名',
  consolePort: '本程序的控制端口',
  normalTcp: '普通TCP推送',
  jtwTcp: '交委TCP推送',
};
function createConfigItem(key, value) {
    const $item = $('<div class="config-item">');
    const displayKey = key.split('.').pop();
    const valueType = typeof value;
    
    let inputElement;
    const isCryptoPacket = displayKey === 'cryptoPacket' || (typeof key === 'string' && key.endsWith('.cryptoPacket'));
    if (isCryptoPacket) {
        let textValue = '';
        if (Array.isArray(value)) {
            textValue = value.map(v => String(v)).join(', ');
        } else if (typeof value === 'string') {
            textValue = value;
        } else if (value != null) {
            textValue = String(value);
        }
        inputElement = `<input type="text" class="config-input" data-key="${key}" value="${textValue}">`;
    } else if (valueType === 'boolean') {
        inputElement = `<select class="config-input" data-key="${key}">
            <option value="true" ${value ? 'selected' : ''}>true</option>
            <option value="false" ${!value ? 'selected' : ''}>false</option>
        </select>`;
    } else if (valueType === 'number') {
        inputElement = `<input type="number" class="config-input" data-key="${key}" value="${value}">`;
    } else {
        inputElement = `<input type="text" class="config-input" data-key="${key}" value="${value}">`;
    }
    
    $item.html(`
        <label>${DISPLAY_LABELS[displayKey] || displayKey}:</label>
        ${inputElement}
    `);
    
    return $item;
}

// 保存配置（支持环境切换时同步迁移顶级节点）
function saveConfig() {
    const updatedConfig = {};
    
    $('.config-input').each(function() {
        const key = $(this).data('key');
        let value = $(this).val();
        
        let converted = value;
        // 特殊处理：cryptoPacket 以逗号分隔保存为数组
        if (typeof key === 'string' && (key.endsWith('.cryptoPacket') || key === 'cryptoPacket')) {
            if (typeof value === 'string') {
                converted = value.split(',').map(s => s.trim()).filter(s => s.length > 0);
            } else if (Array.isArray(value)) {
                converted = value.map(String);
            } else if (value == null) {
                converted = [];
            }
        } else {
            // 类型转换（数字/布尔）
            const originalType = typeof originalConfig[key];
            if (originalType === 'number') {
                converted = parseFloat(value);
            } else if (originalType === 'boolean') {
                converted = value === 'true';
            }
        }
        
        updatedConfig[key] = converted;
    });
    
    // 环境切换：如果 env 发生变化，则将所有以旧环境为前缀的键迁移到新环境前缀
    const oldEnv = originalConfig['env'];
    const newEnv = updatedConfig['env'];
    let finalFlat = updatedConfig;
    if (typeof oldEnv === 'string' && typeof newEnv === 'string' && oldEnv && newEnv && oldEnv !== newEnv) {
        const remapped = {};
        Object.keys(updatedConfig).forEach(k => {
            if (k === 'env') {
                remapped[k] = newEnv;
            } else if (k.startsWith(oldEnv + '.')) {
                remapped[newEnv + '.' + k.slice(oldEnv.length + 1)] = updatedConfig[k];
            } else {
                remapped[k] = updatedConfig[k];
            }
        });
        finalFlat = remapped;
    }
    
    // 将扁平结构转换为嵌套，后端更好地写入 TOML
    const nestedConfig = unflattenToNested(finalFlat);
    
    $.ajax({
        url: `${API_BASE}/setting/save`,
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            config: nestedConfig,
            operation: 'update'
        }),
        success: function(data) {
            if (data.success) {
                originalConfig = JSON.parse(JSON.stringify(finalFlat));
                modifiedKeys.clear();
                $('.config-input, .config-item label').removeClass('modified');
                showConfigStatus('配置保存成功，修改已生效', 'success');
                updateEnvLabels();
            } else {
                showConfigStatus('保存配置失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，保存失败', 'danger');
            console.error('Save config error:', error);
        }
    });
}

// 重新加载配置
function reloadConfig() {
    if (modifiedKeys.size > 0) {
        if (!confirm('当前有未保存的修改，确定要重新加载吗？')) {
            return;
        }
    }
    
    modifiedKeys.clear();
    loadCurrentConfig();
}

// 重置配置
function resetConfig() {
    if (!confirm('确定要重置为默认配置吗？此操作不可撤销！')) {
        return;
    }
    
    $.ajax({
        url: `${API_BASE}/setting/reset`,
        method: 'POST',
        success: function(data) {
            if (data.success) {
                loadCurrentConfig();
                showConfigStatus('配置已重置为默认值', 'warning');
            } else {
                showConfigStatus('重置失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，重置失败', 'danger');
            console.error('Reset config error:', error);
        }
    });
}

// 删除转换连接组（仅针对 env.converter.xxx）
function deleteConfigGroup(groupName) {
    if (!confirm(`确定要删除转换连接 "${groupName}" 吗？`)) {
        return;
    }

    $.ajax({
        url: `${API_BASE}/setting/delete`,
        method: 'DELETE',
        contentType: 'application/json',
        data: JSON.stringify({
            key: groupName,
            operation: 'delete'
        }),
        success: function(data) {
            if (data.success) {
                // 本地移除该组下的所有键
                Object.keys(currentConfig).forEach(function(k) {
                    if (k.startsWith(groupName + '.')) {
                        delete currentConfig[k];
                        delete originalConfig[k];
                        modifiedKeys.delete(k);
                    }
                });
                renderConfigSections();
                showConfigStatus(`转换连接 "${groupName}" 已删除`, 'warning');
            } else {
                showConfigStatus('删除转换连接失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，删除失败', 'danger');
            console.error('Delete group error:', error);
        }
    });
}

// 删除配置项
function deleteConfigItem(key) {
    if (!confirm(`确定要删除配置项 "${key}" 吗？`)) {
        return;
    }
    
    $.ajax({
        url: `${API_BASE}/setting/delete`,
        method: 'DELETE',
        contentType: 'application/json',
        data: JSON.stringify({
            key: key,
            operation: 'delete'
        }),
        success: function(data) {
            if (data.success) {
                delete currentConfig[key];
                delete originalConfig[key];
                modifiedKeys.delete(key);
                renderConfigSections();
                showConfigStatus(`配置项 "${key}" 已删除`, 'warning');
            } else {
                showConfigStatus('删除配置项失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，删除失败', 'danger');
            console.error('Delete config error:', error);
        }
    });
}

// 添加新配置（动态环境）
function addNewConfig() {
    const projectName = $('#new-project-name').val().trim();
    
    if (!projectName) {
        alert('请输入转换连接名称');
        return;
    }
    
    // 验证项目名称格式（只允许字母、数字、下划线）
    if (!/^[a-zA-Z0-9_]+$/.test(projectName)) {
        alert('转换连接名称只能包含字母、数字和下划线');
        return;
    }
    
    const envName = getCurrentEnvironment();
    // 检查是否已存在
    const baseKey = `${envName}.converter.${projectName}`;
    const existingKeys = Object.keys(currentConfig).filter(key => key.startsWith(baseKey + '.'));
    if (existingKeys.length > 0) {
        alert(`转换连接 "${projectName}" 已存在`);
        return;
    }
    
    // 创建转换连接完整 NodeFor809 配置（嵌套结构）
    const subProjectConfigNested = {
        [envName]: {
            converter: {
                [projectName]: {
                    name: projectName,
                    enable: true,
                    cryptoPacket: [],
                    encryptKey: 223344,
                    extendVersion: true,
                    govServerIP: "192.168.3.56",
                    govServerPort: 19001,
                    jtw809ConverterDownLinkIp: "127.0.0.1",
                    jtw809ConverterDownLinkPort: 1302,
                    jtw809ConverterIp: "127.0.0.1",
                    jtw809ConverterPort: 1311,
                    localServerIP: "localhost",
                    localServerPort: 1301,
                    openCrypto: false,
                    platformId: 1001,
                    platformPassword: "123456",
                    platformUserId: 100101,
                    protocolVersion: "1.0.0",
                    thirdpartPort: 11223,
                    useLocationInterval: false
                }
            }
        }
    };
    
    $.ajax({
        url: `${API_BASE}/setting/save`,
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            config: subProjectConfigNested,
            operation: 'add_subproject'
        }),
        success: function(data) {
            if (data.success) {
                // 更新本地配置（扁平化）
                const flatAdded = flattenObject(subProjectConfigNested);
                Object.assign(currentConfig, flatAdded);
                Object.assign(originalConfig, flatAdded);
                
                renderConfigSections();
                clearAddForm();
                showConfigStatus(`转换连接 "${projectName}" 已添加`, 'success');
                updateEnvLabels();
                switchTab('current-config');
            } else {
                alert('添加转换连接失败: ' + data.message);
            }
        },
        error: function(xhr, status, error) {
            alert('网络错误，添加失败');
            console.error('Add subproject error:', error);
        }
    });
}

// 清空添加表单
function clearAddForm() {
    $('#new-project-name').val('');
}

// 加载配置树
function loadConfigTree() {
    const $tree = $('#config-tree');
    $tree.empty();
    
    // 构建配置树结构
    const tree = {};
    Object.keys(currentConfig).forEach(key => {
        const parts = key.split('.');
        let current = tree;
        
        parts.forEach((part, index) => {
            if (!current[part]) {
                current[part] = index === parts.length - 1 ? currentConfig[key] : {};
            }
            current = current[part];
        });
    });
    
    // 渲染树结构
    function renderTree(obj, level = 0) {
        const $ul = $('<ul>');
        
        Object.keys(obj).forEach(key => {
            const $li = $('<li>');
            const indent = '&nbsp;'.repeat(level * 4);
            
            if (typeof obj[key] === 'object' && obj[key] !== null) {
                $li.html(`${indent}<strong>${key}/</strong>`);
                $li.append(renderTree(obj[key], level + 1));
            } else {
                $li.html(`${indent}${key}: <code>${obj[key]}</code>`);
            }
            
            $ul.append($li);
        });
        
        return $ul;
    }
    
    $tree.append(renderTree(tree));
}

// 加载历史记录
function loadHistory() {
    $.ajax({
        url: `${API_BASE}/setting/history`,
        method: 'GET',
        success: function(data) {
            if (data.success) {
                renderHistory(data.history);
            } else {
                $('#history-list').html('<p>加载历史记录失败</p>');
            }
        },
        error: function(xhr, status, error) {
            $('#history-list').html('<p>网络错误，无法加载历史记录</p>');
            console.error('Load history error:', error);
        }
    });
}

// 渲染历史记录
function renderHistory(history) {
    const $list = $('#history-list');
    $list.empty();
    
    if (!history || history.length === 0) {
        $list.html('<p style="text-align: center; padding: 20px; color: #666;">暂无历史记录</p>');
        return;
    }
    
    history.forEach(item => {
        const op = item.operation || {};
        const opType = typeof op.operation === 'string' ? op.operation : (typeof item.operation === 'string' ? item.operation : 'unknown');
        const converterNodes = (op.config && op.config.develop && op.config.develop.converter)
            || (item.full_config && item.full_config.develop && item.full_config.develop.converter)
            || {};
        const converterCount = converterNodes && typeof converterNodes === 'object' ? Object.keys(converterNodes).length : 0;

        const $item = $('<div class="history-item" data-expanded="false">');
        $item.html(`
            <div class="history-summary" style="display:flex;align-items:center;gap:12px;flex-wrap:wrap;">
                <div class="history-time">${item.timestamp}</div>
                <div class="history-op">操作: ${opType}</div>
                <div class="history-count">converter转换连接: ${converterCount}组</div>
                <button class="btn btn-secondary btn-sm toggle-details">详细信息</button>
                <button class="btn btn-secondary btn-sm" onclick="rollbackConfig('${item.timestamp}')">回滚</button>
            </div>
            <div class="history-details" style="display:none;margin-top:8px;">
                <pre style="max-height:300px;overflow:auto;"></pre>
            </div>
        `);
        // 设置JSON文本
        $item.find('pre').text(JSON.stringify(item, null, 2));
        $list.append($item);
    });

    // 展开/收起绑定（使用事件委托避免重复绑定）
    $list.off('click', '.toggle-details').on('click', '.toggle-details', function() {
        const $item = $(this).closest('.history-item');
        const expanded = $item.attr('data-expanded') === 'true';
        $item.attr('data-expanded', expanded ? 'false' : 'true');
        const $details = $item.find('.history-details');
        if (expanded) {
            $details.slideUp(150);
            $(this).text('详细信息');
        } else {
            $details.slideDown(150);
            $(this).text('收起');
        }
    });
}

// 回滚配置
function rollbackConfig(timestamp) {
    if (!confirm(`确定要回滚到 ${timestamp} 的配置吗？`)) {
        return;
    }
    
    $.ajax({
        url: `${API_BASE}/setting/rollback`,
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            timestamp: timestamp,
            operation: 'rollback'
        }),
        success: function(data) {
            if (data.success) {
                loadCurrentConfig();
                showConfigStatus(`已回滚到 ${timestamp} 的配置`, 'warning');
            } else {
                showConfigStatus('回滚失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，回滚失败', 'danger');
            console.error('Rollback config error:', error);
        }
    });
}

// 清空历史
function clearHistory() {
    if (!confirm('确定要清空历史记录吗？此操作不可撤销！')) {
        return;
    }
    
    $.ajax({
        url: `${API_BASE}/setting/clear_history`,
        method: 'POST',
        success: function(data) {
            if (data.success) {
                $('#history-list').empty();
                showConfigStatus('已清空历史记录', 'warning');
            } else {
                showConfigStatus('清空失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，清空失败', 'danger');
            console.error('Clear history error:', error);
        }
    });
}

// 显示配置状态
function showConfigStatus(message, type) {
    const $status = $('#config-status');
    $status.removeClass('alert-success alert-danger alert-warning');
    $status.addClass(`alert-${type}`);
    
    let indicator = 'status-running';
    if (type === 'danger') {
        indicator = 'status-error';
    } else if (type === 'warning') {
        indicator = 'status-modified';
    }
    
    $status.find('.status-indicator').removeClass('status-running status-error status-modified').addClass(indicator);
    $status.find('span').nextAll().remove();
    $status.append(message);
    $status.show();
    
    // 3秒后自动隐藏错误和警告消息
    if (type !== 'success') {
        setTimeout(() => {
            $status.hide();
        }, 3000);
    }
}

// 更新配置状态
function updateConfigStatus() {
    if (modifiedKeys.size > 0) {
        showConfigStatus(`有 ${modifiedKeys.size} 项配置已修改，未重启时显示 * 标识`, 'warning');
    } else {
        showConfigStatus('配置已加载，系统运行正常', 'success');
    }
}

// 更新环境相关文案标签
function updateEnvLabels() {
    const envName = getCurrentEnvironment();
    $('#env-label, #env-label-2').text(envName);
}

// 恢复并增强：加载当前配置（含环境标签更新）
function loadCurrentConfig() {
    $.ajax({
        url: `${API_BASE}/setting/current`,
        method: 'GET',
        success: function(data) {
            if (data.success) {
                const cfg = data.config || {};
                // 兼容嵌套结构：统一扁平化为 dot-keys 以复用现有渲染
                currentConfig = flattenObject(cfg);
                originalConfig = JSON.parse(JSON.stringify(currentConfig));
                renderConfigSections();
                updateEnvLabels();
                showConfigStatus('配置已加载，系统运行正常', 'success');
            } else {
                showConfigStatus('加载配置失败: ' + data.message, 'danger');
            }
        },
        error: function(xhr, status, error) {
            showConfigStatus('网络错误，无法加载配置', 'danger');
            console.error('Load config error:', error);
        }
    });
}

// 标签页切换（缺失补回）
function switchTab(tabId) {
    $('.tab').removeClass('active');
    $(`.tab[data-tab="${tabId}"]`).addClass('active');

    $('.tab-content').removeClass('active');
    $(`#${tabId}`).addClass('active');

    if (tabId === 'history') {
        loadHistory();
    } else if (tabId === 'add-config') {
        loadConfigTree();
    }
}