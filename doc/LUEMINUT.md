# Mexer

Mexer (tai Mekseri suomalaisittain) tai mex-tester on työkalu, joka tarkistaa Abitti\*-kokeiden tenttipakettien eheyden ja aitouden sekä varmistaa, että annetut purkukoodit ovat kelvollisia.

> \*Abitti on Ylioppilastutkintolautakunnan rekisteröimä EU-tavaramerkki (015833742, 015838915).

## Käyttö

Hanki zip-tiedosto, joka sisältää kaikki kokeet, jotka haluat tarkistaa, ja syötä niiden purkukoodit yhteen tekstitiedostoon, eroteltuina uusilla riveillä.

```bash
# Windows
mexer_amd64.exe <zip_tiedoston_polku> <koodi_tiedoston_polku>
# Linux
./mexer_amd64 <zip_tiedoston_polku> <koodi_tiedoston_polku>
```

## Esimerkit

Esimerkki purkukoodi `codes.txt`:

```txt
annostaa syvyytys toukokuu panettaa
kittaus labiili kaatua asettelu
```

Esimerkkikomento:

```bash
# zip-tiedosto exams.zip ja koodit tiedostossa codes.txt

mexer_amd64 ./exams.zip ./codes.txt
```
