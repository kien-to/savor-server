-- Insert sample stores
INSERT INTO stores (
    id, title, description, pickup_time, price, original_price, discounted_price,
    background_url, avatar_url, image_url, rating, reviews, address, items_left,
    latitude, longitude
) VALUES 
    (4, 'Homeskillet Redwood City', 'Surprise Bag', 'Pick up tomorrow 1:00 AM - 5:00 AM',
     5.99, 21.00, 5.99,
     'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
     'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
     'https://images.crowdspring.com/blog/wp-content/uploads/2023/05/16174534/bakery-hero.png',
     4.6, 123, '123 Main St, Redwood City, CA', 5, 
     37.485215, -122.236355),
    (5, 'Pho 75', 'Surprise Bag', 'Pick up tomorrow 1:00 AM - 5:00 AM',
     5.99, 18.00, 5.99,
     'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg',
     'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg',
     'https://www.simplyrecipes.com/thmb/J7YRLoUK0In-BzbTzS1IhFdh_TE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-07-ChickenPho-13-87ae826d1cb347c1a68d133edc7d9a1b.jpg',
     4.6, 89, '456 Broadway St, Redwood City, CA', 3,
     37.486731, -122.232982),
    (6, 'Halal Guys', 'Surprise Bag', 'Pick up tomorrow 7:00 AM - 8:00 AM',
     3.99, 15.00,
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
     'https://tb-static.uber.com/prod/image-proc/processed_images/4f64073782a7b78dadf1605c4c51734b/30be7d11a3ed6f6183354d1933fbb6c7.jpeg',
     4.3, 156, '789 El Camino Real, Palo Alto, CA', 8,
     37.444731, -122.163982);


