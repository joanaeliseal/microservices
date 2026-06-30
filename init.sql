-- Criação dos bancos de dados
CREATE DATABASE IF NOT EXISTS `order_db`;
CREATE DATABASE IF NOT EXISTS `payment_db`;
CREATE DATABASE IF NOT EXISTS `shipping_db`;

-- Usar banco de dados order para inserir itens de estoque
USE `order_db`;

-- Criar tabela de estoque (GORM cria automaticamente, mas podemos popular aqui)
CREATE TABLE IF NOT EXISTS `stock_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_code` varchar(191) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_stock_items_product_code` (`product_code`),
  KEY `idx_stock_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Inserir produtos no estoque para testes
INSERT INTO `stock_items` (`created_at`, `updated_at`, `product_code`, `name`, `quantity`) VALUES
  (NOW(), NOW(), 'PROD-001', 'Notebook', 50),
  (NOW(), NOW(), 'PROD-002', 'Mouse', 200),
  (NOW(), NOW(), 'PROD-003', 'Teclado', 150);
