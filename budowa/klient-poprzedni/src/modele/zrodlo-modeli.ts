import { Command, type Account } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Wykaz identyfikatorów modeli, dla których wolno adresować ustawienia
 * i tożsamość na osi `model`.
 *
 * Kontrakt nie ma komendy `model.list`: model nie jest bytem rejestrowanym,
 * tylko wartością danych, którą niesie kanał modelu (`Channel.model`) albo
 * konto (`Account.defaultModel`). Wykaz składa się z tych dwóch źródeł i służy
 * wyłącznie jako podpowiedź do pola tekstowego — model jeszcze nieużywany
 * wpisuje się identyfikatorem wprost, bo pole nie jest listą zamkniętą.
 *
 * Rejestr kanałów, który nie dotarł, daje podpowiedź pustą, a nie pusty
 * formularz: adresowanie osi modelu działa dalej.
 */
export interface ZrodloModeli {
  /** Identyfikatory modeli znane rdzeniowi, bez powtórzeń, w porządku nazw. */
  identyfikatory(konta: readonly Account[]): Promise<string[]>;
}

export function utworzZrodloModeli(kanal: Kanal): ZrodloModeli {
  return {
    async identyfikatory(konta) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ChannelList, {}),
        Command.ChannelList,
        (tresc) => czyTablica(tresc.channels),
      );

      if (!wynik.udany) {
        console.warn('[modele] rejestr kanałów nie dotarł', wynik.blad?.message ?? '');
      }

      const zKanalow = (wynik.wynik?.channels ?? []).map((kanalModelu) => kanalModelu.model);
      const zKont = konta.map((konto) => konto.defaultModel);
      return bezPowtorzen([...zKanalow, ...zKont]);
    },
  };
}

/** Wykaz bez powtórzeń i bez pozycji pustych, uporządkowany po polsku. */
export function bezPowtorzen(pozycje: readonly (string | undefined)[]): string[] {
  const zebrane = new Set<string>();
  for (const pozycja of pozycje) {
    const oczyszczona = (pozycja ?? '').trim();
    if (oczyszczona !== '') zebrane.add(oczyszczona);
  }
  return [...zebrane].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl'));
}
