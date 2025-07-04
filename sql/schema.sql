-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS chatbot CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- 使用数据库
USE coze;

-- 用户表
CREATE TABLE IF NOT EXISTS user (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户Id',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '密码',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    avatar VARCHAR(255) COMMENT '头像URL',
    role VARCHAR(20) DEFAULT 'user' COMMENT '角色：user/admin',
    status INT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 聊天会话表
CREATE TABLE IF NOT EXISTS conversation (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '会话Id',
    coze_conversation_id VARCHAR(100) NOT NULL COMMENT 'Coze会话Id',
    user_id INT UNSIGNED NOT NULL COMMENT '用户Id',
    title VARCHAR(100) COMMENT '会话标题',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 消息表
CREATE TABLE IF NOT EXISTS message (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '消息Id',
    coze_message_id VARCHAR(100) NOT NULL COMMENT 'Coze消息Id',
    conversation_id INT UNSIGNED NOT NULL COMMENT '会话Id',
    model_id INT UNSIGNED NOT NULL COMMENT 'AI模型Id',
    metadata VARCHAR(255) COMMENT '元数据，如模型参数等',
    role VARCHAR(10) NOT NULL COMMENT 'user或assistant',
    content TEXT NOT NULL COMMENT '消息内容',
    tokens INT DEFAULT 0 COMMENT '消耗Token数量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_chat_id (conversation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
