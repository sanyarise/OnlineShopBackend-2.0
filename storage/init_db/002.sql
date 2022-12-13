CREATE EXTENSION pgcrypto;

CREATE TABLE rights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256),
    rules text[]
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL,
    lastname VARCHAR(256),
    password TEXT NOT NULL,
    email VARCHAR(256) NOT NULL UNIQUE,
    rights UUID,
    zipcode VARCHAR(16),
    country VARCHAR(256),
    city VARCHAR(256),
    street VARCHAR(256),
    CONSTRAINT fk_rights 
        FOREIGN KEY(rights) REFERENCES rights(id)
);


CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL UNIQUE, 
    description TEXT NOT NULL,
    picture TEXT
);


CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL,
    category UUID,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    vendor TEXT NOT NULL,
    pictures text[],
    CONSTRAINT fk_category
        FOREIGN KEY(category) REFERENCES categories(id)
);

CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expire_at timestamp NOT NULL DEFAULT now() + interval '1 hour',
    user_id UUID,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE cart_items (
    cart_id UUID,
    item_id UUID,
    PRIMARY KEY(cart_id, item_id),
    CONSTRAINT fk_cart_id
        FOREIGN KEY(cart_id) REFERENCES carts(id),
    CONSTRAINT fk_item_id
        FOREIGN KEY(item_id) REFERENCES items(id)    
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_time timestamp not NULL,
    user_id UUID,
    status VARCHAR(256),
    address TEXT,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id) REFERENCES users(id)
);


CREATE TABLE order_items (
    order_id UUID,
    item_id UUID,
    PRIMARY KEY(order_id, item_id),
    CONSTRAINT fk_order_id
        FOREIGN KEY(order_id) REFERENCES orders(id),
    CONSTRAINT fk_item_id
        FOREIGN KEY(item_id) REFERENCES items(id)
);

INSERT INTO categories(id, name, description, picture) VALUES ('d0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'electronics', 'electronics for life', 'http://localhost:8000/files/categories/d0d3df2d-f6c8-4956-9d76-998ee1ec8a39/20221213125935.jpeg');
INSERT INTO categories(id, name, description, picture) VALUES ('ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'clothes', 'clothes for life', 'http://localhost:8000/files/categories/ad1c1a7f-3210-4554-b0d7-6c57d4e04657/20221213130104.jpeg');
INSERT INTO categories(id, name, description, picture) VALUES ('57804957-c219-4874-8962-d8a08eb368da', 'music', 'music for life', 'http://localhost:8000/files/categories/57804957-c219-4874-8962-d8a08eb368da/20221213130339.jpeg');

INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('0b74b0ac-68aa-462b-8609-4bf5eac3f9f7', 'smartphone samsung', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 10000, 'samsung', '{"http://localhost:8000/files/items/0b74b0ac-68aa-462b-8609-4bf5eac3f9f7/20221213132612.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('692a759d-a993-45ea-b3bd-6cd523db74b4', 'smartphone xiaomi', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 8900, 'xiaomi', '{"http://localhost:8000/files/items/692a759d-a993-45ea-b3bd-6cd523db74b4/20221213132650.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('d9674ae5-ed88-4e0f-a80d-2388186c211b', 'smartphone vivo', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 9900, 'vivo', '{"http://localhost:8000/files/items/d9674ae5-ed88-4e0f-a80d-2388186c211b/20221213132703.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('6395c53c-437c-4cb4-968c-fa4c7d14d32a', 'smartphone vivo2', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 19990, 'vivo', '{"http://localhost:8000/files/items/6395c53c-437c-4cb4-968c-fa4c7d14d32a/20221213132716.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('e0ada75f-389c-493f-953b-7f3ee48b5a00', 'smartphone vivo3', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 14990, 'samsung', '{"http://localhost:8000/files/items/e0ada75f-389c-493f-953b-7f3ee48b5a00/20221213132728.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('2c0ba661-ffb3-4085-8b1c-0303c92eacaa', 'smartphone bq', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 4990, 'bq', '{"http://localhost:8000/files/items/2c0ba661-ffb3-4085-8b1c-0303c92eacaa/20221213132743.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('13be90f4-377b-4de8-8dc2-3fcfb5706129', 'smartphone bq', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 3990, 'samsung', '{"http://localhost:8000/files/items/13be90f4-377b-4de8-8dc2-3fcfb5706129/20221213132756.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('223ba075-8c96-4760-8620-f0adbe3f2e7e', 'smartphone bq', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 6990, 'samsung', '{"http://localhost:8000/files/items/223ba075-8c96-4760-8620-f0adbe3f2e7e/20221213132809.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('3f05a859-1806-4851-a457-c9cb69b22846', 'smartphone techno', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 9990, 'samsung', '{"http://localhost:8000/files/items/3f05a859-1806-4851-a457-c9cb69b22846/20221213132821.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('7eb10fc6-20e8-4b3b-95a3-0f55cf0145d6', 'smartphone techno', 'd0d3df2d-f6c8-4956-9d76-998ee1ec8a39', 'best smartphone', 8990, 'samsung', '{"http://localhost:8000/files/items/7eb10fc6-20e8-4b3b-95a3-0f55cf0145d6/20221213132833.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('f15e6c99-e47a-49ff-a8f7-4a8efc77c946', 'skirt', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best skirt', 1190, 'adidas', '{"http://localhost:8000/files/items/f15e6c99-e47a-49ff-a8f7-4a8efc77c946/20221213132844.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('a9914762-59a0-4faa-b278-4236c39a5eda', 'skirt2', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best skirt', 1290, 'adidas', '{"http://localhost:8000/files/items/a9914762-59a0-4faa-b278-4236c39a5eda/20221213132857.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('407054c4-8efc-47bd-b298-2eb1d2ef74d2', 'skirt3', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best skirt', 1690, 'adidas', '{"http://localhost:8000/files/items/407054c4-8efc-47bd-b298-2eb1d2ef74d2/20221213132911.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('a847a063-c078-4a8e-ac7a-62356db88d16', 'skirt4', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best skirt', 2290, 'zolla', '{"http://localhost:8000/files/items/a847a063-c078-4a8e-ac7a-62356db88d16/20221213132937.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('97af7789-10e6-4f9c-8f14-8295d00294cf', 'skirt5', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best skirt', 1390, 'zolla', '{"http://localhost:8000/files/items/97af7789-10e6-4f9c-8f14-8295d00294cf/20221213132949.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('2b3c9ac6-54f5-4217-acde-6f5dab670aae', 't-short', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best t-short', 600, 'gloria jeans', '{"http://localhost:8000/files/items/2b3c9ac6-54f5-4217-acde-6f5dab670aae/20221213133044.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('b682b27c-669b-46e9-b2a9-1b42ac01fd01', 't-short2', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best t-short', 700, 'adidas', '{"http://localhost:8000/files/items/b682b27c-669b-46e9-b2a9-1b42ac01fd01/20221213133057.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('3d0ecdf6-2062-4cbb-a267-1c90eb87d51b', 'shoes', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best shoes', 4990, 'nike', '{"http://localhost:8000/files/items/3d0ecdf6-2062-4cbb-a267-1c90eb87d51b/20221213133110.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('584b0269-bc2f-49bb-a7aa-757790b6c843', 'shoes', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best shoes', 5990, 'adidas', '{"http://localhost:8000/files/items/584b0269-bc2f-49bb-a7aa-757790b6c843/20221213133129.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('20855daf-07bd-42e5-8901-e8beecad40dc', 't-short', 'ad1c1a7f-3210-4554-b0d7-6c57d4e04657', 'best t-short', 390, 'adidas', '{"http://localhost:8000/files/items/20855daf-07bd-42e5-8901-e8beecad40dc/20221213134220.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('981564ff-9cec-491d-98ed-93534b6b4ac5', 'mp3 player', '57804957-c219-4874-8962-d8a08eb368da', 'best mp3', 1990, 'sony', '{"http://localhost:8000/files/items/981564ff-9cec-491d-98ed-93534b6b4ac5/20221213133150.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('18b6847d-0657-43f8-b0aa-ab08d4095ae2', 'mp3 player2', '57804957-c219-4874-8962-d8a08eb368da', 'best mp3', 2990, 'sony', '{"http://localhost:8000/files/items/18b6847d-0657-43f8-b0aa-ab08d4095ae2/20221213133204.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('8fc13002-dab1-4df3-bb26-6e089f3e949e', 'mp3 player3', '57804957-c219-4874-8962-d8a08eb368da', 'best mp3', 990, 'sony', '{"http://localhost:8000/files/items/8fc13002-dab1-4df3-bb26-6e089f3e949e/20221213133218.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('c687c4ca-e670-4e00-a5fd-e4f49f3763db', 'boombox', '57804957-c219-4874-8962-d8a08eb368da', 'best boombox', 5000, 'BBK', '{"http://localhost:8000/files/items/c687c4ca-e670-4e00-a5fd-e4f49f3763db/20221213133232.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('bcb04300-b292-4cb3-82b3-35228abd83a4', 'boombox2', '57804957-c219-4874-8962-d8a08eb368da', 'best boombox', 4990, 'BBK', '{"http://localhost:8000/files/items/bcb04300-b292-4cb3-82b3-35228abd83a4/20221213133247.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('fe182012-8d07-4205-8f2c-3fe106271e9f', 'boombox3', '57804957-c219-4874-8962-d8a08eb368da', 'best boombox', 8990, 'sony', '{"http://localhost:8000/files/items/fe182012-8d07-4205-8f2c-3fe106271e9f/20221213133303.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('d135cc36-858a-4380-87e9-7833c728291a', 'earphones', '57804957-c219-4874-8962-d8a08eb368da', 'best earphones', 1000, 'smartbuy', '{"http://localhost:8000/files/items/d135cc36-858a-4380-87e9-7833c728291a/20221213133320.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('7c60376b-4373-448b-a859-1b3a43e9b47b', 'earphones2', '57804957-c219-4874-8962-d8a08eb368da', 'best earphones', 1990, 'sony', '{"http://localhost:8000/files/items/7c60376b-4373-448b-a859-1b3a43e9b47b/20221213133333.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('bb5d9cce-5759-476d-8fd0-01af7a6af84a', 'earphones3', '57804957-c219-4874-8962-d8a08eb368da', 'best earphones', 2990, 'huawei', '{"http://localhost:8000/files/items/bb5d9cce-5759-476d-8fd0-01af7a6af84a/20221213133347.jpeg"}');
INSERT INTO items(id, name, category, description, price, vendor, pictures) VALUES ('642a21ec-3c7b-40bc-bbf5-15179dfb2ddb', 'earphones4', '57804957-c219-4874-8962-d8a08eb368da', 'best earphones', 5990, 'samsung', '{"http://localhost:8000/files/items/642a21ec-3c7b-40bc-bbf5-15179dfb2ddb/20221213133401.jpeg"}');



