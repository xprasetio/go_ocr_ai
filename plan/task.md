- Buatkan apps sistem OCR menggunakan : 
BE : golang echo, (DDD + DI Sarulabs + RawSql + Viper + Cobra + Clean Code)
FE : NextJs + tailwind
Database: Cloudflare D1 (SQL). Wajib data harus persist di D1

Untuk struktur folder sesuaikan dengan stucture_folder.md agar memudahkan dalam membuat project gunakan schema_db untuk databasenya sertakan juga generate database saat applikasi dijalankan sertakan validasi : jika sudah pernah di generate maka abaikan

Untuk Design sesuaikan dengan design.md dan asset.html sebagai acuan utama

GEMINI_MODEL=gemini-2.5-flash
GEMINI_API_KEY=AIzaSyDjNSfMPJ
hoZZ28foKaXdYCFQi_umUgLrU
PORT=8081

Syarat utama pembuatan yaitu menggunakan model gemini dengan api_key
kemudian buatkan : 

file_attachment berupa png atau pdf, membaca dokumen keuangan (resi/struk/invoice) tambahkan prompt didalamnya : dengan prompt contoh sebagai berikut : 
analisa gambar yang diinputkan, solusikan juga kasus dokumen buruk: miring, gelap, blur, atau bukan resi, jika ada angka 0 dibelakang titik kurang dari tiga maka hapus, contoh 10000.0 menjadi 10000 kemudian buatkan response output seperti berikut 

{
  "amount": 0,
  "change": 0,
  "created_at": "string",
  "date_shopping": "string",
  "id": 0,
  "items": [
    {
      "created_at": "string",
      "expenses_id": 0,
      "id": 0,
      "name": "string",
      "price": 0,
      "qty": 0,
      "subtotal": 0,
      "updated_at": "string"
    }
  ],
  "note": "string",
  "parsed_data": "string",
  "receipt_image": "string",
  "title": "string",
  "updated_at": "string"
}
- Ketika tombol anlyize ditekan makan nanti akan muncul preview disamping kanan dan editable dikarenakan OCR tidak akurat 100%
- untuk bucket file storage disimpan di supabase dengan nama : ocr_ai_receipt
- Buatkan list Recent Transactions dengan pagination bukan dengan card
- Pencarian berdasarkan tanggal dan nama vendor
- ada tombol detail
- export csv 
- semua yang bersifat rahasia ditaruh di .env seperti GEMINI_API_KEY
