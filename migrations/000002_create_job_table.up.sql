CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    number_of_bedrooms TEXT NOT NULL,
    additional_services TEXT,
    description_additional_services TEXT,
    truck_size TEXT NOT NULL,
    pickup_datetime TIMESTAMP NOT NULL,
    delivery_datetime TIMESTAMP NOT NULL,
    cut_amount DOUBLE PRECISION,
    payment_amount DOUBLE PRECISION
);