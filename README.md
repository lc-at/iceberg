# Iceberg
Iceberg is a WhatsApp chatbot designed to manage school assignments. This is my first Golang project. It was built using go-whatsapp library (https://github.com/Rhymen/go-whatsapp). This project is not well-documented but some things should be self-explanatory.

## Installation
Before installing, make sure you have a working MySQL database, a device with WhatsApp installed, and a good internet connection.
```
go get github.com/p4kl0nc4t/iceberg
go install github.com/p4kl0nc4t/iceberg
```
## Usage
1. Set up a configuration file based on ![this](https://github.com/p4kl0nc4t/iceberg/blob/master/config.yml.default) template.
2. Run `iceberg -config path_to_config.yml`.
3. Scan the shown barcode.
4. It is up and running!

## License
This project is licensed with MIT License.

## Contribution
Feel free to contribute to this project. Any kind of contribution is really appreciated.
