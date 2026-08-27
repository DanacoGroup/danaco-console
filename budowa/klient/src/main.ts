/**
 * Złożenie klienta — jedyny punkt wejścia pakietu.
 *
 * Wyłącznie kompozycja: gniazdo, kanał, sesja, tożsamość, przebieg, montaż.
 * Zero rozgałęzień drogi wejścia i zero dotknięcia dokumentu poza wskazaniem
 * korzenia montażowi — etap i odsłonę rozstrzyga `przebieg.ts`, a węzły
 * stawia `montaz.ts`.
 *
 * Magazyn tokenu bramki zostaje domyślny, czyli w pamięci procesu. Magazyn
 * trwały należy do powłoki i żadne źródło go dziś nie wskazuje; magazyn
 * udający trwałość obiecywałby rozpoznanie urządzenia, którego nie ma.
 */

import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from './polaczenie/adres-rdzenia.ts';
import { utworzTransport } from './polaczenie/gniazdo.ts';
import { utworzKanal } from './protokol/kanal.ts';
import { utworzSesje } from './protokol/sesja.ts';
import { tozsamoscKlienta } from './protokol/tozsamosc-klienta.ts';
import { zamontuj } from './wejscie/montaz.ts';
import { utworzPrzebieg } from './wejscie/przebieg.ts';

/**
 * Adres gniazda rdzenia.
 *
 * Rdzeń wystawia pakiet interfejsu obok gniazda (`uchwytStatyki`
 * w `server/internal/transport/statyka.go`), więc dokument wczytany po HTTP
 * przyszedł z rdzenia i to jego adres jest adresem gniazda. Dokument wczytany
 * inaczej — z pliku albo z protokołu powłoki — pochodzenia nie niesie, więc
 * zostaje pętla zwrotna z portem domyślnym.
 */
function adresRdzenia(): string {
  return adresGniazdaRdzenia(globalThis.location.origin) ?? adresRdzeniaLokalnego();
}

const klient = tozsamoscKlienta();
const transport = utworzTransport(adresRdzenia());
const przebieg = utworzPrzebieg({
  kanal: utworzKanal(transport, utworzSesje()),
  transport,
  klient,
});

zamontuj({ korzen: document, przebieg, wersjaKlienta: klient.wersja });
przebieg.polacz();
