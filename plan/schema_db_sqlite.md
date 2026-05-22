CREATE TABLE expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    date_shopping TEXT NULL,
    "change" INTEGER NULL,
    amount REAL NULL,
    note TEXT NULL,
    receipt_image TEXT NULL,
    parsed_data TEXT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE expense_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expenses_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    qty INTEGER NOT NULL,
    price INTEGER NOT NULL,
    subtotal INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (expenses_id)
        REFERENCES expenses(id)
        ON DELETE CASCADE
);