/**
 * Plik stanowi jedyny punkt wejścia klienta i składa gniazdo, kanał, sesję,
 * tożsamość, przebieg oraz montaż węzłów, nie wprowadzając własnych
 * rozgałęzień drogi wejścia.
 */

import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from './polaczenie/adres-rdzenia.ts';
import { utworzTransport } from './polaczenie/gniazdo.ts';
import { utworzKanal } from './protokol/kanal.ts';
import { utworzSesje } from './protokol/sesja.ts';
import { tozsamoscKlienta } from './protokol/tozsamosc-klienta.ts';
import { zamontuj } from './wejscie/montaz.ts';
import { utworzPrzebieg } from './wejscie/przebieg.ts';

/**
 * Adres gniazda rdzenia pochodzi z dokumentu wczytanego po HTTP; w pozostałych
 * przypadkach zostaje pętla zwrotna z portem domyślnym.
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
