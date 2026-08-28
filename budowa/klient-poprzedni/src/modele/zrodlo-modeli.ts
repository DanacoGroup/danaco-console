import { Command, type Account } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Wykaz identyfikatorów modeli, dla których wolno adresować ustawienia
 * i tożsamość na osi `model`. Wykaz składa się z kanałów rdzenia oraz z kont
 * i służy jako podpowiedź do pola tekstowego, które nie jest listą zamkniętą.
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

/**
 * Wykaz bez powtórzeń i bez pozycji pustych, uporządkowany porównaniem napisów
 * właściwym dla polszczyzny; pozycje puste oraz złożone z samych odstępów
 * odpadają.
 */
export function bezPowtorzen(pozycje: readonly (string | undefined)[]): string[] {
  const zebrane = new Set<string>();
  for (const pozycja of pozycje) {
    const oczyszczona = (pozycja ?? '').trim();
    if (oczyszczona !== '') zebrane.add(oczyszczona);
  }
  return [...zebrane].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl'));
}
