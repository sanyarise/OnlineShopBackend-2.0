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