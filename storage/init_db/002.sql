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
    picture TEXT,
    deleted_at timestamptz NULL
);


CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL,
    category UUID,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    vendor TEXT NOT NULL,
    pictures text[],
    deleted_at timestamptz NULL,
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
    item_quantity INTEGER NOT NULL,
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

INSERT INTO rights(id) VALUES ('c8bfa6d6-38c7-458e-8e37-c98940713d16');

INSERT INTO users(id, name, password, email, rights) VALUES ('ad9e481d-1064-4430-89af-adc6f9a5751c', 'anonymus', '0000', 'example@example.com', 'c8bfa6d6-38c7-458e-8e37-c98940713d16');























