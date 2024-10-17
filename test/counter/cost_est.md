**Total Products**:  
Rakuten has approximately **31.5 million** products.

**Data Sources**:  
The product data can be accessed via the following URLs:

-   [https://books.rakuten.co.jp/json/genre/000](https://books.rakuten.co.jp/json/genre/000)
-   [https://books.rakuten.co.jp/json/genre/001](https://books.rakuten.co.jp/json/genre/001)
-   [https://books.rakuten.co.jp/json/genre/002](https://books.rakuten.co.jp/json/genre/002)
-   [https://books.rakuten.co.jp/json/genre/003](https://books.rakuten.co.jp/json/genre/003)
-   [https://books.rakuten.co.jp/json/genre/004](https://books.rakuten.co.jp/json/genre/004)
-   [https://books.rakuten.co.jp/json/genre/005](https://books.rakuten.co.jp/json/genre/005)
-   [https://books.rakuten.co.jp/json/genre/006](https://books.rakuten.co.jp/json/genre/006)
-   [https://books.rakuten.co.jp/json/genre/007](https://books.rakuten.co.jp/json/genre/007)
-   [https://books.rakuten.co.jp/json/genre/101](https://books.rakuten.co.jp/json/genre/101)

**Current Crawling Speed**:

-   **Rate**: ~**2000 products/hour** per VM
-   **Configuration**: 100 concurrent requests (8 vCPUs & 12 GB RAM)

**Monthly Crawling Capacity**:

-   Products crawled per month: 2000 products/hour×730 hours=1,460,000 products2000 \text{ products/hour} \times 730 \text{ hours} = 1,460,000 \text{ products}2000 products/hour×730 hours=1,460,000 products

**Target Completion**:  
To complete crawling **31.5 million** products within **30 days**, we can consider two approaches:

### **Approaches**

1.  **Using Multiple VMs**:

    -   **Required VMs**: **22**
    -   **Total Products Crawled**:

    22 VMs×1,460,000 products/VM=32,120,000 products22 \text{ VMs} \times 1,460,000 \text{ products/VM} = 32,120,000 \text{ products}22 VMs×1,460,000 products/VM=32,120,000 products
    -   **Estimated Cost**: **$4,594**
2.  **Using a Single VM with Higher Resources**:

    -   **Required vCPUs**: **176** (22 VMs * 8 vCPUs)
    -   **Required RAM**: **264 GB** (22 VMs * 12 GB)
    -   **Estimated Cost**: **$5,723.95**