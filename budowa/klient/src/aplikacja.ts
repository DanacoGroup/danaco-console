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
import { utworzPrzebieg, type StanPrzebiegu } from './wejscie/przebieg.ts';
import { zamontujRame } from './rama/montaz.ts';
import { gotowaDoPrzekazania } from './rama/przekazanie.ts';

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

const oknoWejscia = zamontuj({ korzen: document, przebieg, wersjaKlienta: klient.wersja });

/* Scena drogi wejścia i miejsce ramy stoją w dokumencie jako dwa węzły
   odrębne od tego, co montuje `zamontuj` wewnątrz sceny — przekazanie
   przełącza między nimi raz, nieodwracalnie. */
let przekazano = false;

function naZmianePrzebiegu(stan: StanPrzebiegu): void {
  if (przekazano || !gotowaDoPrzekazania(stan)) return;
  przekazano = true;

  const scenaWejscia = document.querySelector('[data-wejscie]');
  const miejsceRamy = document.querySelector('[data-rama-aplikacji]');
  if (scenaWejscia === null || miejsceRamy === null || stan.srodowisko === undefined) return;

  oknoWejscia.zdejmij();
  (scenaWejscia as HTMLElement).hidden = true;
  zamontujRame({
    miejsce: miejsceRamy as HTMLElement,
    srodowisko: stan.srodowisko,
    moduly: stan.moduly,
    sesje: stan.sesje,
  });
  (miejsceRamy as HTMLElement).hidden = false;
}

przebieg.naZmiane(naZmianePrzebiegu);
przebieg.polacz();
