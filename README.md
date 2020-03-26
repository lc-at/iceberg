# Iceberg
Iceberg is a WhatsApp chatbot designed to manage school assignments. This is my first Golang project. It was built using go-whatsapp library (https://github.com/Rhymen/go-whatsapp). It also uses a MySQL database to save assignment records. This project is not well-documented but some things should be self-explanatory.

## Installation
If you don't want to build from source, go ahead and check out the ![Releases](https://github.com/ttycelery/iceberg/releases) page.
### Building from source
```
git clone https://github.com/ttycelery/iceberg
cd iceberg
go build .
```
## Usage
Before using, make sure you have a working MySQL database, a device with WhatsApp installed, and a good internet connection.
1. Set up a configuration file based on ![this](https://github.com/ttycelery/iceberg/blob/master/config.yml.default) template.
2. Run `iceberg -config path_to_config.yml`.
3. Scan the shown barcode.
4. It is up and running!

## License
This project is licensed with MIT License.

## Contribution
Feel free to contribute to this project. Any kind of contribution is really appreciated.
