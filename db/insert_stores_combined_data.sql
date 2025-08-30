-- Create stores table and insert data from stores_combined_data.csv

-- Drop table if it exists (optional - uncomment if you want to recreate the table)
-- DROP TABLE IF EXISTS stores CASCADE;

-- Create the stores table
CREATE TABLE IF NOT EXISTS stores (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    pickup_time VARCHAR(255),
    distance VARCHAR(50),
    price DECIMAL(10,2),
    original_price DECIMAL(10,2),
    discounted_price DECIMAL(10,2),
    background_url TEXT,
    avatar_url TEXT,
    rating DECIMAL(3,1),
    reviews INTEGER,
    address TEXT,
    items_left INTEGER,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    google_maps_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    image_url TEXT
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_stores_location ON stores(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_stores_title ON stores(title);
CREATE INDEX IF NOT EXISTS idx_stores_description ON stores(description);
CREATE INDEX IF NOT EXISTS idx_stores_address ON stores(address);

-- Clear existing data (optional - uncomment if you want to start fresh)
-- DELETE FROM stores;

-- Insert Vietnamese restaurant data (Hanoi)
INSERT INTO stores (id, title, description, pickup_time, price, original_price, discounted_price, background_url, avatar_url, image_url, rating, reviews, address, items_left, latitude, longitude) VALUES
('1', 'Pho Gia Truyen', 'Traditional Vietnamese pho and bun bo hue', 'Pick up today 5:00 PM - 8:00 PM', 8.50, 18.00, 8.50, 'https://dynamic-media-cdn.tripadvisor.com/media/photo-o/0e/d4/f0/1f/omg.jpg?w=900&h=500&s=1', 'https://dynamic-media-cdn.tripadvisor.com/media/photo-o/0e/d4/f0/1f/omg.jpg?w=900&h=500&s=1', 'https://dynamic-media-cdn.tripadvisor.com/media/photo-o/0e/d4/f0/1f/omg.jpg?w=900&h=500&s=1', 4.7, 456, '49 Bat Dan Hoan Kiem Hanoi', 6, 21.0347, -105.8563),

('2', 'Bun Cha Huong Lien', 'Famous bun cha where Obama dined', 'Pick up tomorrow 11:00 AM - 2:00 PM', 12.00, 25.00, 12.00, 'https://i0.wp.com/hungryghostfoodandtravel.com/wp-content/uploads/2024/12/Vietnamese-Chicken-Meatball-Noodle-Bowl_done.jpg?resize=800%2C530&quality=89&ssl=1', 'https://i0.wp.com/hungryghostfoodandtravel.com/wp-content/uploads/2024/12/Vietnamese-Chicken-Meatball-Noodle-Bowl_done.jpg?resize=800%2C530&quality=89&ssl=1', 'https://i0.wp.com/hungryghostfoodandtravel.com/wp-content/uploads/2024/12/Vietnamese-Chicken-Meatball-Noodle-Bowl_done.jpg?resize=800%2C530&quality=89&ssl=1', 4.8, 892, '24 Le Van Huu Hai Ba Trung Hanoi', 4, 21.0186, -105.8492),

('3', 'Banh Mi 25', 'Crispy Vietnamese baguette sandwiches', 'Pick up today 2:00 PM - 6:00 PM', 3.50, 8.00, 3.50, 'https://silverkris.singaporeair.com/wp-content/uploads/2022/08/Banh-Mi-25.jpg', 'https://silverkris.singaporeair.com/wp-content/uploads/2022/08/Banh-Mi-25.jpg', 'https://silverkris.singaporeair.com/wp-content/uploads/2022/08/Banh-Mi-25.jpg', 4.5, 234, '25 Hang Ca Hoan Kiem Hanoi', 8, 21.0285, -105.8542),

('4', 'Quan An Ngon', 'Traditional Vietnamese street food', 'Pick up tomorrow 6:00 PM - 9:00 PM', 15.00, 30.00, 15.00, 'https://quanngonrestaurant.com/wp-content/uploads/2023/09/Quan-Ngon-Banner.jpg', 'https://quanngonrestaurant.com/wp-content/uploads/2023/09/Quan-Ngon-Banner.jpg', 'https://quanngonrestaurant.com/wp-content/uploads/2023/09/Quan-Ngon-Banner.jpg', 4.6, 678, '18 Phan Boi Chau Hoan Kiem Hanoi', 5, 21.0192, -105.8530),

('5', 'Cha Ca La Vong', 'Famous turmeric fish with dill and noodles', 'Pick up tomorrow 7:00 PM - 9:00 PM', 18.00, 35.00, 18.00, 'https://feedthepudge.com/wp-content/uploads/2025/03/Cha-Ca-La-Vong-Cover-Photo.webp', 'https://feedthepudge.com/wp-content/uploads/2025/03/Cha-Ca-La-Vong-Cover-Photo.webp', 'https://feedthepudge.com/wp-content/uploads/2025/03/Cha-Ca-La-Vong-Cover-Photo.webp', 4.4, 345, '14 Cha Ca La Vong Hanoi', 3, 21.0344, -105.8520),

('6', 'Bun Bo Nam Bo', 'Signature Hanoi beef noodle salad', 'Pick up tomorrow 12:00 PM - 3:00 PM', 9.00, 16.00, 9.00, 'https://cdn.sanity.io/images/jdiyrv6o/production/cf566e8e131a49721247b77566dc1cf17237aa1e-700x467.jpg?auto=format&w=3840', 'https://cdn.sanity.io/images/jdiyrv6o/production/cf566e8e131a49721247b77566dc1cf17237aa1e-700x467.jpg?auto=format&w=3840', 'https://cdn.sanity.io/images/jdiyrv6o/production/cf566e8e131a49721247b77566dc1cf17237aa1e-700x467.jpg?auto=format&w=3840', 4.6, 289, '67 Hang Dieu Hoan Kiem Hanoi', 7, 21.0310, -105.8545),

('7', 'Xoi Yen', 'Traditional sticky rice with toppings', 'Pick up today 7:00 AM - 10:00 AM', 5.50, 12.00, 5.50, 'https://dynamic-media-cdn.tripadvisor.com/media/photo-o/02/2d/b2/63/xoi-yen.jpg?w=900&h=500&s=1', 'https://dynamic-media-cdn.tripadvisor.com/media/photo-o/02/2d/b2/63/xoi-yen.jpg?w=900&h=500&s=1', 'https://dynamic-media-cdn.tripadvisor.com/media/photo-o/02/2d/b2/63/xoi-yen.jpg?w=900&h=500&s=1', 4.3, 167, '35B Nguyen Huu Huan Hoan Kiem Hanoi', 9, 21.0289, -105.8548),

-- Insert San Francisco restaurant data
('sf001', 'Ferry Building Marketplace', 'Artisanal breads and pastries from local bakery', 'Pick up tomorrow 9:00 PM - 10:00 PM', 12.99, 25.99, 12.99, 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 4.8, 234, '1 Ferry Building San Francisco', 8, 37.7955, -122.3937),

('sf002', 'Chinatown Bakery', 'Traditional Chinese baked goods and dim sum', 'Pick up tomorrow 8:30 PM - 9:30 PM', 8.99, 18.99, 8.99, 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 4.6, 189, '123 Grant Ave San Francisco', 5, 37.7927, -122.4067),

('sf003', 'Mission District Cafe', 'Fresh Mexican pastries and coffee', 'Pick up tomorrow 7:00 PM - 8:00 PM', 10.99, 22.99, 10.99, 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 4.7, 189, '2128 Mission St San Francisco', 6, 37.7599, -122.4148),

('sf004', 'Hayes Valley Patisserie', 'French pastries and artisanal coffee', 'Pick up tomorrow 8:00 PM - 9:00 PM', 15.99, 32.99, 15.99, 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 4.9, 278, '432 Octavia St San Francisco', 4, 37.7759, -122.4245),

-- Insert Bay Area restaurant data (Redwood City, Palo Alto)
('rc001', 'Homeskillet Redwood City', 'Surprise Bag', 'Pick up tomorrow 1:00 AM - 5:00 AM', 5.99, 21.00, 5.99, 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png', 4.6, 123, '123 Main St Redwood City CA', 5, 37.485215, -122.236355),

('rc002', 'Pho 75', 'Surprise Bag', 'Pick up tomorrow 1:00 AM - 5:00 AM', 5.99, 18.00, 5.99, 'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg', 'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg', 'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg', 4.6, 89, '456 Broadway St Redwood City CA', 3, 37.486731, -122.232982),

('pa001', 'Halal Guys', 'Surprise Bag', 'Pick up tomorrow 7:00 AM - 8:00 AM', 3.99, 15.00, 3.99, 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg', 4.3, 156, '789 El Camino Real Palo Alto CA', 8, 37.444731, -122.163982);

-- Verify the insertion
SELECT COUNT(*) as total_stores FROM stores;
SELECT id, title, address FROM stores ORDER BY id;
