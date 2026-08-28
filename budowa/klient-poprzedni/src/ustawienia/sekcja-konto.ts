import { Command } from '../../../shared/contract';
import { utworzSekcjeBraku } from './sekcja-braku';
import type { SekcjaUstawien } from './sekcje';

/** Sekcja konta operatora nazywa miejsce bez ani jednego pola, bo kontrakt nie niesie dziś komendy odczytu profilu operatora. */
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
