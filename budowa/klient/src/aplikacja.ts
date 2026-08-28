/**
 * Jedyny punkt wejścia klienta. Składa gniazdo, kanał, sesję, tożsamość
 * i przebieg drogi wejścia; gdy przebieg dochodzi do środowiska, oddaje
 * sterowanie ramie aplikacji i zdejmuje scenę drogi wejścia z ekranu.
 */

import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from './polaczenie/adres-rdzenia.ts';
import { utworzTransport } from './polaczenie/gniazdo.ts';
import { utworzKanal } from './protokol/kanal.ts';
import { utworzSesje } from './protokol/sesja.ts';
import { tozsamoscKlienta } from './protokol/tozsamosc-klienta.ts';
import { zamontuj } from './wejscie/montaz.ts';
import { utworzPrzebieg } from './wejscie/przebieg.ts';
import { utworzPrzekazanieJednorazowe } from './rama/przekazanie.ts';

/**
 * Adres gniazda rdzenia pochodzi z dokumentu wczytanego po HTTP; w pozostałych
 * przypadkach zostaje pętla zwrotna z portem domyślnym.
 */
function adresRdzenia(): string {
  return adresGniazdaRdzenia(globalThis.location.origin) ?? adresRdzeniaLokalnego();
}

const klient = tozsamoscKlienta();
const transport = utworzTransport(adresRdzenia());
const kanal = utworzKanal(transport, utworzSesje());
const przebieg = utworzPrzebieg({ kanal, transport, klient });

const oknoWejscia = zamontuj({ korzen: document, przebieg, wersjaKlienta: klient.wersja });

/* Scena drogi wejścia i miejsce ramy stoją w dokumencie jako dwa węzły
   odrębne od tego, co montuje `zamontuj` wewnątrz sceny — przekazanie przełącza
   między nimi raz, dopiero po udanym przekazaniu, więc odmowa (np. węzeł
   montażu jeszcze nie stoi w dokumencie) nie blokuje próby przy kolejnej
   zmianie stanu. */
const naZmianePrzebiegu = utworzPrzekazanieJednorazowe({
  dokument: document,
  zdejmijOknoWejscia: () => oknoWejscia.zdejmij(),
  kanal,
});

przebieg.naZmiane(naZmianePrzebiegu);
przebieg.polacz();
