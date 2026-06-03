# SB Wallet Backend

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Project for gin exercise by Koda batch

## Technologies Used

- [![Gin-Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-green?logo=Gin&logoColor=white)](https://gin-gonic.com/en/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.4-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)
- etc.

## Features

- Authentication (Register, Login, Logout)
- Dashboard (Balance, History, )
- etc.

## Usage Instruction

### DB setup

1. Create your environment on the root directory named `.env`

```
DB_HOST={YOUR_DB_HOST}
DB_PORT={YOUR_DB_PORT}
DB_PASS={YOUR_DB_PASS}
DB_NAME={YOUR_DB_NAME}
DB_USER={YOUR_DB_USER}
```

### Running the Application

1. Clone this repository

```bash
$ git clone url
```

2. Install dependency

```bash
$ go mod download
```

3. etc.

## Routes

### Features

| Endpoint  | Method | Description |
| --------- | ------ | ----------- |
| /auth     | POST   | Login       |
| /auth/new | POST   | Register    |

### Documentation

For complete documentation, visit `/swagger/index.html`

## Changelog

| Version | Desctiption        |
| ------- | ------------------ |
| 1.0.0   | Initial app        |
| 1.1.0   | add authentication |

- 1.1.1
  - Fixed bug in .... by [L1mus](https://github.com/L1mus)
  - Fixed bug in .... by [AnggaVb](https://github.com/anggavb/)
- 1.2.0
  - Add .... feature by [BernadDwiki](https://github.com/BernadDwiki)
  - Add .... feature by [rivando-al-rasyid](https://github.com/rivando-al-rasyid/)

## How to Contribute

- Fork this repository
- Create your changes
- Pull Request

## License

This project is licensed under the MIT License

<!-- ## CONTACTS
[email](mailto:) -->

## Related Project

[Frontend](https://github.com/nopalllfd/koda-b7-react)
