// Etap trzeci drogi wejścia: przygotowanie środowiska pracy. Wykaz etapów
// i pasek postępu idą za odpowiedziami rdzenia.
import { Command, SessionStatus } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { zapomnijTokenSesji } from '../protokol/token-sesji.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { napis, urzadzenieTrwale } from './wejscie-katalog.ts';
import { zdanieOdmowy } from './wejscie-odmowa.ts';

const NAGLOWEK = 'Wejście';

const ETAP_UWIERZYTELNIENIA = 0;
const ETAP_SESJI = 2;

const ODMIANY: Readonly<Record<string, string>> = {
  gotowy: 'dn-krok--poprawny',
  pracuje: 'dn-krok--pracuje',
  blad: 'dn-krok--wstrzymany',
};

const ZNAK_GOTOWY = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" '
  + 'stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">'
  + '<path d="M20 6 9 17l-5-5"/></svg>';

let domkniete = false;
let domknij: () => void = () => undefined;

/* Etapy bez pokrycia w kontrakcie zostają w stanie oczekiwania: wykaz mówi
   wtedy, czego droga wejścia nie robi, zamiast malować rzecz zrobioną. */
export async function przygotujSrodowisko(kanal: Kanal, poGotowosci: () => void): Promise<void> {
  domkniete = false;
  domknij = () => {
    if (domkniete) return;
    domkniete = true;
    poGotowosci();
  };
  const kroki = wyzerujEtapy();
  const czeka = napis('przygotowanie.miary.oczekuje');

  oznacz(kroki[ETAP_UWIERZYTELNIENIA], 'pracuje', czeka);
  await rozpoznajUrzadzenie(kanal, kroki[ETAP_UWIERZYTELNIENIA]);
  ustawPostep(zrobione(kroki));

  oznacz(kroki[ETAP_SESJI], 'pracuje', czeka);
  await policzSesje(kanal, kroki[ETAP_SESJI]);
  ustawPostep(zrobione(kroki));
  domknij();
}

export function wyzerujEtapy(): HTMLElement[] {
  ustawPostep(0);
  const kroki = wykazKrokow();
  const czeka = napis('przygotowanie.miary.oczekuje');
  for (const krok of kroki) oznacz(krok, 'oczekuje', czeka);
  return kroki;
}

export function zwiazPasPrzygotowania(kanal: Kanal): void {
  wyzerujEtapy();
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('.we-pas [data-idz="uwierzytelnienie"]') !== null) {
      void wyloguj(kanal);
      return;
    }
    if (cel.closest('.we-pas [data-komunikat-tytul]') !== null) domknij();
    if (cel.closest('[data-zmien-haslo]') !== null) void zmienHaslo(kanal);
  }, true);
}

// Zmiana hasła ze znanym hasłem bieżącym; odzyskanie listem to osobna auth.reset.
async function zmienHaslo(kanal: Kanal): Promise<void> {
  const biezace = document.querySelector<HTMLInputElement>('[data-haslo-biezace]');
  const nowe = document.querySelector<HTMLInputElement>('[data-haslo-nowe]');
  if (biezace === null || nowe === null) return;
  const wynik = await wywolaj(kanal, Command.AuthPasswordReset, {
    currentPassword: biezace.value,
    newPassword: nowe.value,
  });
  if (!wynik.udany) {
    oglos('Konto', wynik.blad?.message ?? 'Rdzeń odmówił zmiany hasła.', 'ostrzezenie');
    return;
  }
  biezace.value = '';
  nowe.value = '';
  oglos('Konto', 'Hasło bramki zmienione.');
}

/* Wylogowanie unieważnia sesję bramki tego urządzenia — kontrakt nie zna innej
   drogi jej zamknięcia. Token schodzi niezależnie od odpowiedzi rdzenia. */
async function wyloguj(kanal: Kanal): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.DeviceList, {});
  const biezace = wykaz.wynik?.devices.find((urzadzenie) => urzadzenie.current);
  const wskazane = biezace?.deviceId ?? urzadzenieTrwale();
  const zdjecie = await wywolaj(kanal, Command.DeviceRevoke, { deviceId: wskazane });
  zapomnijTokenSesji();
  if (!zdjecie.udany) {
    oglos(NAGLOWEK, zdanieOdmowy(zdjecie.blad), 'blad');
    return;
  }
  oglos(NAGLOWEK, `Sesja bramki urządzenia ${wskazane} została unieważniona.`);
}

async function rozpoznajUrzadzenie(kanal: Kanal, krok: HTMLElement | undefined): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeviceList, {});
  const czeka = napis('przygotowanie.miary.oczekuje');
  if (!wynik.udany || wynik.wynik === undefined) {
    oznacz(krok, 'blad', czeka);
    return;
  }
  const biezace = wynik.wynik.devices.find((urzadzenie) => urzadzenie.current);
  if (biezace === undefined) {
    oznacz(krok, 'oczekuje', czeka);
    return;
  }
  oznacz(krok, 'gotowy', napis('przygotowanie.miary.rozpoznane'));
}

async function policzSesje(kanal: Kanal, krok: HTMLElement | undefined): Promise<void> {
  const wynik = await wywolaj(kanal, Command.SessionList, { includePresence: true });
  if (!wynik.udany || wynik.wynik === undefined) {
    oznacz(krok, 'blad', napis('przygotowanie.miary.oczekuje'));
    return;
  }
  const czynne = wynik.wynik.sessions.filter(
    (sesja) => sesja.status === SessionStatus.Active,
  ).length;
  oznacz(krok, 'gotowy', podstaw(napis('przygotowanie.miary.karty'), {
    odtworzone: String(czynne),
    wszystkie: String(wynik.wynik.total),
  }));
}

function wykazKrokow(): HTMLElement[] {
  const okno = document.querySelector<HTMLElement>('[data-wejscie-okno="przygotowanie"]');
  return [...(okno?.querySelectorAll<HTMLElement>('.dn-krok') ?? [])];
}

function zrobione(kroki: HTMLElement[]): number {
  if (kroki.length === 0) return 0;
  const gotowe = kroki.filter((krok) => krok.dataset.stan === 'gotowy').length;
  return Math.round((gotowe / kroki.length) * 100);
}

function oznacz(krok: HTMLElement | undefined, stan: string, miara: string): void {
  if (krok === undefined) return;
  krok.dataset.stan = stan;
  for (const [nazwa, odmiana] of Object.entries(ODMIANY)) {
    krok.classList.toggle(odmiana, nazwa === stan);
  }
  const meta = krok.querySelector('.dn-krok-meta');
  if (meta !== null) meta.textContent = miara;
  const znak = krok.querySelector('.dn-krok-znak');
  if (znak === null) return;
  if (stan === 'gotowy') znak.innerHTML = ZNAK_GOTOWY;
  else if (stan === 'pracuje') znak.innerHTML = '<span class="dn-kropka dn-kropka--tetno"></span>';
}

function ustawPostep(procent: number): void {
  const okno = document.querySelector<HTMLElement>('[data-wejscie-okno="przygotowanie"]');
  if (okno === null) return;
  const liczba = okno.querySelector('.pg-postep-wiersz b');
  if (liczba !== null) liczba.textContent = `${procent}%`;
  okno.querySelector('.dn-postep-tor')?.setAttribute('aria-valuenow', String(procent));
  const wartosc = okno.querySelector<HTMLElement>('.dn-postep-wartosc');
  if (wartosc === null) return;
  wartosc.dataset.wartosc = String(procent);
  wartosc.style.width = `${procent}%`;
}

function podstaw(wzor: string, dane: Record<string, string>): string {
  return wzor.replace(/\{(\w+)\}/g, (calosc, nazwa: string) => dane[nazwa] ?? calosc);
}
