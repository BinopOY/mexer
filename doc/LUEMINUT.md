# Mexer

Mexer (tai Mekseri suomalaisittain) tai mex-tester on työkalu, joka tarkistaa Abitti\*-kokeiden koepakettien eheyden ja aitouden sekä varmistaa, että annetut purkukoodit ovat kelvollisia.

> \*Abitti on Ylioppilastutkintolautakunnan rekisteröimä EU-tavaramerkki (015833742, 015838915).

## Asentaminen

Voit ladata viimeisimmän version windowsille tai linuxille osoitteesta https://github.com/BinopOY/mexer/releases/latest

## Käyttö

Mikä on koepaketti:

-   Zip-tiedosto, joka sisältää useita kokeita.
-   Zip-tiedosto, joka sisältää yhden kokeen.
-   Mex-tiedosto yhtenä kokeena.

Mikä on purkukoodin syöte:

-   Yksittäinen purkukoodi, joka annetaan komentoriviargumenttina.
-   . \*.txt-tiedosto, joka sisältää luetttelon purkukoodeista, eroteltuna rivinvaihdolla.

Purkukoodeja tulee olla vähintään yhtä monta kuin tarkistettavia kokeita.

```bash
# Windows
mexer_amd64.exe <koepaketin_polku> <purkukoodi>
# Linux
./mexer_amd64 <koepaketin_polku> <purkukoodi>
```

## Esimerkit

Esimerkki purkukooditiedosto `codes.txt`:

```txt
annostaa syvyytys toukokuu panettaa
kittaus labiili kaatua asettelu
```

Esimerkkikomentoja:

```bash
# zip-paketti tiedostossa exams.zip ja koodit tiedostossa codes.txt
mexer_amd64 ./exams.zip ./codes.txt

# Yksi koe zip-pakettina ja purkukoodi syötteenä
mexer_amd64 ./exams.zip "kittaus labiili kaatua asettelu"

# Yksi koe mex-tiedostona ja purkukoodi syötteenä
mexer_amd64 ./exams.mex "kittaus labiili kaatua asettelu"

# Yksi koe mex-tiedostona ja purkukoodit tiedostossa codes.txt
mexer_amd64 ./exams.mex ./codes.txt
```
