CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  order_uid TEXT UNIQUE NOT NULL,
  track_number TEXT,
  entry TEXT,
  locale TEXT,
  internal_signature TEXT,
  customer_id TEXT,
  delivery_service TEXT,
  shard_key TEXT,
  sm_id INTEGER,
  date_created TEXT,
  oof_shard TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  
  transaction TEXT NOT NULL,
  request_id TEXT NOT NULL,
  currency TEXT NOT NULL,
  provider TEXT NOT NULL,
  amount INTEGER NOT NULL,
  payment_dt INTEGER NOT NULL,
  bank TEXT NOT NULL,
  delivery_cost INTEGER NOT NULL,
  goods_total INTEGER NOT NULL,
  custom_fee INTEGER NOT NULL,

  name TEXT NOT NULL,
  phone TEXT NOT NULL,
  zip TEXT NOT NULL,
  city TEXT NOT NULL,
  address TEXT NOT NULL,
  region TEXT NOT NULL,
  email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  order_uid TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id INTEGER NOT NULL,
  track_number TEXT NOT NULL,
  price INTEGER NOT NULL,
  rid TEXT NOT NULL,
  name TEXT NOT NULL,
  sale INTEGER NOT NULL,
  size TEXT NOT NULL,
  total_price INTEGER NOT NULL,
  nm_id INTEGER NOT NULL,
  brand TEXT NOT NULL,
  status INTEGER NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_orders_order_uid ON orders(order_uid);


