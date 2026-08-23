import { Command } from '../../../shared/contract';
import { utworzSekcjeBraku } from './sekcja-braku';
import type { SekcjaUstawien } from './sekcje';

/**
 * Sekcja „Konto Operatora" — miejsce nazwane, bez ani jednego pola.
 *
 * Konto ISTNIEJE: rejestracja zakłada jedyną encję właściciela z loginem,
 * adresem e-mail uwierzytelniającym, stanem potwierdzenia i datą utworzenia.
 * Sekcja nie ma jednak czym go pokazać, bo brakuje KOMENDY ODCZYTU profilu —
 * `auth.register` konto zapisuje, `auth.verify` oddaje sesję, ale żadna komenda
 * nie zwraca loginu ani adresu. Pole bez komendy odczytu byłoby puste, a pole
 * bez komendy zapisu — atrapą, więc ich tu nie ma do czasu, aż rodzina odczytu
 * profilu wejdzie do kontraktu.
 *
 * Rodzina `account.*` nie jest tym, czego tu brakuje: opisuje konta modeli
 * (klucze dostawców), nie konto Operatora, i jest obsłużona
 * w `budowa/client/src/modele/`. Sekcja „Konta modeli" niżej w rejestrze
 * odsyła właśnie tam.
 */
export function utworzSekcjeKonto(): SekcjaUstawien {
  return utworzSekcjeBraku({
    wstep:
      'Konto Operatora istnieje, ale nie ma go dziś czym pokazać w tym oknie: ' +
      'rdzeń zna login, adres e-mail uwierzytelniający, stan potwierdzenia i datę ' +
      'utworzenia, lecz brakuje komendy, która by je odczytała. Pola bez komendy ' +
      'odczytu byłyby puste, więc czekają na rodzinę odczytu profilu w kontrakcie.',
    pomiar: [
      'Plik migracja_125_konto_wlasciciela.sql zakłada encję `konto_wlasciciela` ' +
        '(login, email, potwierdzone, utworzono, warunek `id = 1`) oraz ' +
        '`potwierdzenie_tozsamosci`. Warstwa danych to ' +
        '`server/internal/dane/konto_wlasciciela.go`.',
      `Rejestracja (\`${Command.AuthRegister}\`) konto zakłada, a potwierdzenie adresu ` +
        `(\`${Command.AuthVerify}\`) przenosi je do stanu potwierdzonego — obie komendy je ` +
        'zapisują, żadna nie zwraca jego pól.',
      'Czego brakuje, to KOMENDA ODCZYTU profilu Operatora — kontrakt nie niesie ' +
        'jej dziś w żadnej rodzinie. Dopóki jej nie ma, sekcja nazywa konto, ' +
        'ale go nie wyświetla.',
    ],
    odeslanie:
      'Hasło bramki, PIN urządzenia i sesje zmienia się w sekcji ' +
      '„Uwierzytelnianie". Konta modeli (klucze dostawców) to rodzina account.* ' +
      'i osobna sekcja „Konta modeli".',
  });
}
