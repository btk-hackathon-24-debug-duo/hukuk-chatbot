# Hukuk Chatbot Projesi

Bu proje, hukuk öğrencilerinin ve avukatların kolayca bilgilere ulaşmasını sağlayan bir chatbot uygulamasıdır. Frontend kısmında Next.js, backend kısmında ise Go kullanılmıştır.

## Kurulum

### Frontend

1. Proje dizinine gidin:
    ```bash
    cd apps/frontend
    ```
2. Gerekli paketleri yükleyin:
    ```bash
    npm install
    ```
3. Uygulamayı başlatın:
    ```bash
    npm run dev
    ```

### Backend

1. Proje dizinine gidin:
    ```bash
    cd apps/backend
    ```
2. Gerekli paketleri yükleyin:
    ```bash
    go mod tidy
    ```
3. Uygulamayı başlatın:
    ```bash
    go run main.go
    ```

## Kullanım

1. Frontend ve backend sunucularını başlattıktan sonra, tarayıcınızda `http://localhost:3000` adresine gidin.
2. Chatbot arayüzü üzerinden sorularınızı sorarak bilgi alabilirsiniz.

## Katkıda Bulunma

Katkıda bulunmak isterseniz, lütfen bir pull request gönderin veya bir issue açın.

## Lisans

Bu proje MIT Lisansı ile lisanslanmıştır.