-- Migracja 131 zasila tabelę szablon_studio pięcioma układami fabrycznymi pisma, umowy, raportu, notatki i oferty ze znacznikami pól do wypełnienia.

INSERT INTO szablon_studio (identyfikator_zewnetrzny, nazwa, opis, format, tresc, pola_json, fabryczny)
VALUES
    ('studio-szablon-pismo', 'Pismo', 'Pismo urzędowe z nagłówkiem nadawcy i adresata', 'markdown',
     '{{nadawca}}

{{miejscowosc}}, {{data}}

{{adresat}}

**Dotyczy: {{sprawa}}**

{{tresc}}

Z poważaniem
{{podpis}}',
     '[{"name":"nadawca","label":"Nadawca","required":true},
       {"name":"miejscowosc","label":"Miejscowość","required":true},
       {"name":"data","label":"Data","required":true},
       {"name":"adresat","label":"Adresat","required":true},
       {"name":"sprawa","label":"Czego dotyczy","required":true},
       {"name":"tresc","label":"Treść pisma","required":true},
       {"name":"podpis","label":"Podpis","required":false}]', 1),

    ('studio-szablon-umowa', 'Umowa', 'Umowa dwustronna z paragrafami i miejscem na załączniki', 'markdown',
     '# {{tytul}}

zawarta {{data}} w {{miejscowosc}} pomiędzy:

**{{strona_pierwsza}}**

a

**{{strona_druga}}**

## Przedmiot umowy

{{przedmiot}}

## Termin i wynagrodzenie

{{termin_i_wynagrodzenie}}

## Postanowienia końcowe

{{postanowienia}}

Załączniki: {{zalaczniki}}',
     '[{"name":"tytul","label":"Tytuł umowy","required":true},
       {"name":"data","label":"Data zawarcia","required":true},
       {"name":"miejscowosc","label":"Miejscowość","required":true},
       {"name":"strona_pierwsza","label":"Strona pierwsza","required":true},
       {"name":"strona_druga","label":"Strona druga","required":true},
       {"name":"przedmiot","label":"Przedmiot umowy","required":true},
       {"name":"termin_i_wynagrodzenie","label":"Termin i wynagrodzenie","required":true},
       {"name":"postanowienia","label":"Postanowienia końcowe","required":false},
       {"name":"zalaczniki","label":"Załączniki","required":false}]', 1),

    ('studio-szablon-raport', 'Raport', 'Raport ze streszczeniem zarządczym na wstępie', 'markdown',
     '# {{tytul}}

**Autor:** {{autor}} · **Data:** {{data}}

## Streszczenie zarządcze

{{streszczenie}}

## Ustalenia

{{ustalenia}}

## Wnioski

{{wnioski}}

## Rekomendacje

{{rekomendacje}}',
     '[{"name":"tytul","label":"Tytuł raportu","required":true},
       {"name":"autor","label":"Autor","required":true},
       {"name":"data","label":"Data","required":true},
       {"name":"streszczenie","label":"Streszczenie zarządcze","required":true},
       {"name":"ustalenia","label":"Ustalenia","required":true},
       {"name":"wnioski","label":"Wnioski","required":true},
       {"name":"rekomendacje","label":"Rekomendacje","required":false}]', 1),

    ('studio-szablon-notatka', 'Notatka', 'Notatka ze spotkania z listą ustaleń i działań', 'markdown',
     '# {{temat}}

**Data:** {{data}} · **Uczestnicy:** {{uczestnicy}}

## Przebieg

{{przebieg}}

## Ustalenia

{{ustalenia}}

## Działania

{{dzialania}}',
     '[{"name":"temat","label":"Temat spotkania","required":true},
       {"name":"data","label":"Data","required":true},
       {"name":"uczestnicy","label":"Uczestnicy","required":true},
       {"name":"przebieg","label":"Przebieg","required":false},
       {"name":"ustalenia","label":"Ustalenia","required":true},
       {"name":"dzialania","label":"Działania","required":true}]', 1),

    ('studio-szablon-oferta', 'Oferta', 'Oferta handlowa z zakresem, ceną i terminem ważności', 'markdown',
     '# {{tytul}}

**Dla:** {{odbiorca}} · **Data:** {{data}} · **Ważna do:** {{wazna_do}}

## Zakres

{{zakres}}

## Warunki i cena

{{warunki}}

## Termin realizacji

{{termin}}

{{uwagi}}',
     '[{"name":"tytul","label":"Tytuł oferty","required":true},
       {"name":"odbiorca","label":"Odbiorca","required":true},
       {"name":"data","label":"Data wystawienia","required":true},
       {"name":"wazna_do","label":"Ważna do","required":true},
       {"name":"zakres","label":"Zakres oferty","required":true},
       {"name":"warunki","label":"Warunki i cena","required":true},
       {"name":"termin","label":"Termin realizacji","required":true},
       {"name":"uwagi","label":"Uwagi","required":false}]', 1)
ON CONFLICT(identyfikator_zewnetrzny) DO NOTHING;
