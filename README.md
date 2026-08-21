# Monitor de RTX 5090 para México

Busca tarjetas **RTX 5090** en DDTech y Cyberpuerta. Son comercios mexicanos, por lo que los precios se tratan como MXN y la compra no depende de una importación individual. El monitor excluye tiendas extranjeras: sus impuestos, aranceles y cargos de paquetería sólo se confirman al pagar, así que no se puede prometer que no habrá extras.

## Configuración

1. Copia `.env.example` a un gestor de secretos o configura sus variables en tu sistema. No subas un archivo `.env`.
2. Ajusta `MAX_PRICE_MXN` (por defecto, 80,000).
3. Si usas Gmail, genera una **App Password** con verificación en dos pasos y úsala como `SMTP_PASSWORD`; no uses tu contraseña normal.
4. Ejecuta `go run ./cmd`.

El proyecto no contiene direcciones de correo fijas. Debes definir `EMAIL_FROM`, `EMAIL_TO` y `SMTP_USERNAME`. `PSWRD` se acepta temporalmente por compatibilidad, pero debe migrarse a `SMTP_PASSWORD`.

## GitHub Actions

Para que el monitor programado funcione, crea estos **Actions secrets** en `Settings → Secrets and variables → Actions`: `SMTP_USERNAME`, `SMTP_PASSWORD`, `EMAIL_FROM` y `EMAIL_TO`. Puedes crear como variables no secretas `MAX_PRICE_MXN`, `SMTP_HOST` y `SMTP_PORT` si quieres cambiar los valores predeterminados.

Para validar únicamente Gmail, ejecuta **RTX 5090 monitor → Run workflow**, activa `Send an SMTP test email without checking stores` y lanza el flujo. Recibirás un correo con el asunto `Prueba del monitor RTX 5090`. Esta prueba no consulta tiendas y permite distinguir un problema SMTP de un bloqueo del catálogo.

Cada ejecución normal informa en sus logs cuántas RTX 5090 fueron detectadas y cuántas estaban disponibles por tienda. Si una tienda falla, el monitor continúa con las demás; sólo falla completamente cuando ninguna fuente puede consultarse.

Si alguna vez hubo una clave real en el `.env` que estaba rastreado, rótala primero. Quitar el archivo de la versión actual no elimina las copias de los commits anteriores.

Los precios mostrados son los de catálogo. Antes de comprar, confirma disponibilidad, código postal, costo de envío y factura en el checkout.
