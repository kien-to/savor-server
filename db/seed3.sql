-- First, delete existing data
DELETE FROM store_highlights;
DELETE FROM stores WHERE id IN ('sf001', 'sf002', 'sf003', 'sf004', 'sf005');

-- Add stores in different San Francisco neighborhoods
INSERT INTO stores
    (
    id, title, description, pickup_time, price, original_price,
    background_url, avatar_url, image_url, rating, reviews, address, items_left,
    latitude, longitude
    )
VALUES
    (
        'sf001',
        'Ferry Building Marketplace',
        'Artisanal breads and pastries from local bakery',
        'Pick up tomorrow 9:00 PM - 10:00 PM',
        12.99,
        25.99,
     'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
     'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
     'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
        4.8,
        234,
        '1 Ferry Building, San Francisco',
        8,
        37.7955,
        -122.3937
  ),
    (
        'sf002',
        'Chinatown Bakery',
        'Traditional Chinese baked goods and dim sum',
        'Pick up tomorrow 8:30 PM - 9:30 PM',
        8.99,
        18.99,
   'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg',
     'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg',
     'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg',
        4.6,
        156,
        '759 Clay St, San Francisco',
        5,
        37.7941,
        -122.4078
  ),
    (
        'sf003',
        'Mission District Cafe',
        'Fresh Mexican pastries and coffee',
        'Pick up tomorrow 7:00 PM - 8:00 PM',
        10.99,
        22.99,
        'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
        'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
        4.7,
        189,
        '2128 Mission St, San Francisco',
        6,
        37.7599,
        -122.4148
  ),
    (
        'sf004',
        'Hayes Valley Patisserie',
        'French pastries and artisanal coffee',
        'Pick up tomorrow 8:00 PM - 9:00 PM',
        15.99,
        32.99,
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
        4.9,
        278,
        '432 Octavia St, San Francisco',
        4.0,
        37.7759,
        -122.4245
  ),
    ('sf005', 'Homeskillet Redwood City', 'Surprise Bag', 'Pick up tomorrow 1:00 AM - 5:00 AM',
        5.99, 21.00,
        'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
        'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
        'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
        4.6, 123, '123 Main St, Redwood City, CA', 5,
        37.485215, -122.236355);
