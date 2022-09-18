CREATE EXTENSION pgcrypto;

CREATE TABLE rights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256),
    rules text[]
);

INSERT INTO rights (name, rules) 
    VALUES 
        ('user', ARRAY['place orders']),
        ('admin', ARRAY['delete users', 'add items']);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL,
    password TEXT,
    email VARCHAR(256) NOT NULL UNIQUE,
    rights UUID,
    CONSTRAINT fk_rights 
        FOREIGN KEY(rights) REFERENCES rights(id)
);


CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL UNIQUE, 
    description TEXT NOT NULL
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
    PRIMARY KEY(cart_id, item_id)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_time timestamp not NULL,
    user_id UUID,
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