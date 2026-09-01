// Brama kontraktu sprawdza żądanie wobec kontraktu przed czynnością domeny: pola wymagane i wartości wyliczeń, czytane z artefaktu kontraktu; powitanie kanału jest jedynym wyjątkiem.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"sync"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// znacznikPolaKontraktu nazywa znacznik struktury Go, spod którego brama odczytuje nazwę pola żądania odpowiadającą polu kontraktu.
const znacznikPolaKontraktu = "json"

// wskazaniePolaOpcjonalnego to człon znacznika, którym generator kontraktu oznacza pole opcjonalne; pole wymagane niesie znacznik bez tego członu, wprost z oznaczenia „wymagane” kontraktu.
const wskazaniePolaOpcjonalnego = "omitempty"

// komendaPozaBrama nazywa jedyną komendę, której brama nie sprawdza. Stała
// zamiast warunku wpisanego w gałąź, bo wyjątek ma być odczytywalny w jednym
// miejscu — wyjątek rozsypany po warunkach przestaje być jedynym.
const komendaPozaBrama = shared.CommandConnectionHello

// sprawdzZadanieWobecKontraktu odmawia żądaniu niezgodnemu z kontraktem: sprawdza obecność pól wymaganych oraz przynależność wartości do wyliczenia kontraktu, oba wyczytane z kontraktu.
func sprawdzZadanieWobecKontraktu(ctx context.Context, komenda shared.MessageType,
	ladunek json.RawMessage, wzor any) error {

	pola, sa := polaTresci(ladunek)
	if !sa {
		return nil
	}
	brakujace := brakujacePolaWymagane(wzor, pola)
	if komenda == komendaPozaBrama {
		odnotujBrakiPowitania(ctx, komenda, brakujace)
		return nil
	}
	if len(brakujace) > 0 {
		return bladZgodnosciZKontraktem(komenda,
			"żądanie bez wartości w polach wymaganych kontraktem: "+strings.Join(brakujace, ", "))
	}
	return niezgodnoscWyliczen(komenda, wzor, pola)
}

// odnotujBrakiPowitania kładzie braki powitania w dzienniku rdzenia, jedynym miejscu, do którego mają dojść, ponieważ odpowiedź powitania pola na braki nie niesie.
func odnotujBrakiPowitania(ctx context.Context, komenda shared.MessageType, brakujace []string) {
	if len(brakujace) == 0 {
		return
	}
	dziennik := dziennikZKontekstu(ctx)
	if dziennik == nil {
		return
	}
	dziennik.Printf("brama kontraktu: %s bez pól wymaganych kontraktem: %s; "+
		"powitanie odpowiada mimo to, bo klient odczytuje z niego wersję protokołu",
		komenda, strings.Join(brakujace, ", "))
}

// polaTresci rozkłada treść żądania na pola wierzchnie. Drugi wynik odróżnia treść nierozkładalną na pola od treści rozłożonej, lecz pustej.
func polaTresci(ladunek json.RawMessage) (map[string]json.RawMessage, bool) {
	if len(ladunek) == 0 {
		return map[string]json.RawMessage{}, true
	}
	var pola map[string]json.RawMessage
	if err := json.Unmarshal(ladunek, &pola); err != nil {
		return nil, false
	}
	if pola == nil {
		return map[string]json.RawMessage{}, true
	}
	return pola, true
}

// ksztaltZadania oddaje kształt struktury treści żądania. Wzór przychodzi
// z obsługiwacza komendy, więc bywa wskaźnikiem; typ niebędący strukturą pól
// kontraktu nie niesie i sprawdzać w nim nie ma czego.
func ksztaltZadania(wzor any) reflect.Type {
	ksztalt := reflect.TypeOf(wzor)
	for ksztalt != nil && ksztalt.Kind() == reflect.Pointer {
		ksztalt = ksztalt.Elem()
	}
	if ksztalt == nil || ksztalt.Kind() != reflect.Struct {
		return nil
	}
	return ksztalt
}

// brakujacePolaWymagane zwraca nazwy pól, które kontrakt oznacza jako wymagane,
// a których treść żądania nie niesie albo niesie bez wartości. Nazwy idą
// w kolejności kontraktu, bo w tej kolejności czyta je człowiek w opisie
// komendy.
func brakujacePolaWymagane(wzor any, pola map[string]json.RawMessage) []string {
	ksztalt := ksztaltZadania(wzor)
	if ksztalt == nil {
		return nil
	}
	var brakujace []string
	for i := range ksztalt.NumField() {
		pole := ksztalt.Field(i)
		nazwa, wymagane := polePola(pole)
		if !wymagane {
			continue
		}
		surowa, jest := pola[nazwa]
		if !jest || !poleNiesieWartosc(pole.Type, surowa) {
			brakujace = append(brakujace, nazwa)
		}
	}
	return brakujace
}

// poleNiesieWartosc odróżnia wartość pola od samej obecności klucza. `null`
// nie jest wartością pola napisowego, a pusty napis nie jest wartością pola
// wyliczeniowego, bo żadne wyliczenie kontraktu pustej wartości nie zna —
// przepuszczone dochodzą do zapisu i zatrzymuje je dopiero CHECK sterownika
// bazy, wracając surowym tekstem sterownika zamiast odmową walidacji. Pusty
// napis w polu swobodnym zostaje wartością, bo kontrakt zna komendy, w których
// znaczy zdjęcie wskazania. Pole wielowartościowe rozstrzyga się samą
// obecnością, bo wycinek pusty jest wartością, a Go zapisuje go jako `null`.
func poleNiesieWartosc(typ reflect.Type, surowa json.RawMessage) bool {
	if typ.Kind() != reflect.String {
		return true
	}
	if bytes.Equal(bytes.TrimSpace(surowa), []byte("null")) {
		return false
	}
	var napis string
	if err := json.Unmarshal(surowa, &napis); err != nil {
		return true
	}
	return napis != "" || nazwaWyliczenia(typ) == ""
}

// polePola odczytuje ze znacznika struktury nazwę pola kontraktu oraz to, czy
// pole jest wymagane. Pole nieeksportowane i pole wyłączone znacznikiem „-"
// nie należą do treści kontraktu.
func polePola(pole reflect.StructField) (string, bool) {
	if !pole.IsExported() {
		return "", false
	}
	znacznik, jest := pole.Tag.Lookup(znacznikPolaKontraktu)
	if !jest {
		return "", false
	}
	nazwa, reszta, _ := strings.Cut(znacznik, ",")
	if nazwa == "" || nazwa == "-" {
		return "", false
	}
	for _, czlon := range strings.Split(reszta, ",") {
		if czlon == wskazaniePolaOpcjonalnego {
			return nazwa, false
		}
	}
	return nazwa, true
}

// niezgodnoscWyliczen odmawia żądaniu, którego pole wyliczeniowe niesie wartość spoza zakresu kontraktu; pola pozostałych typów zakresu nie mają i przechodzą.
func niezgodnoscWyliczen(komenda shared.MessageType, wzor any, pola map[string]json.RawMessage) error {
	ksztalt := ksztaltZadania(wzor)
	if ksztalt == nil {
		return nil
	}
	for i := range ksztalt.NumField() {
		pole := ksztalt.Field(i)
		nazwa, _ := polePola(pole)
		if nazwa == "" {
			continue
		}
		wyliczenie := nazwaWyliczenia(pole.Type)
		if wyliczenie == "" {
			continue
		}
		wartosc, jest := pola[nazwa]
		if !jest {
			continue
		}
		zakres, znany := zakresyWyliczenKontraktu()[wyliczenie]
		if !znany {
			return bladWyliczeniaBezZakresu(komenda, nazwa, wyliczenie)
		}
		if powod := wartoscPozaZakresem(nazwa, wartosc, zakres); powod != "" {
			return bladZgodnosciZKontraktem(komenda, powod)
		}
	}
	return nil
}

// nazwaWyliczenia oddaje nazwę wyliczenia kontraktu stojącego pod typem pola,
// albo pusty napis, gdy pole wyliczenia nie niesie. Generator kontraktu
// wystawia wyliczenie nazwanym typem napisowym, pole opcjonalne wskaźnikiem,
// a pole wielowartościowe wycinkiem — stąd zdejmowanie obu powłok. Typ
// napisowy bez własnego pakietu to zwykły `string`, czyli pole bez wyliczenia.
func nazwaWyliczenia(typ reflect.Type) string {
	for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.String || typ.PkgPath() == "" {
		return ""
	}
	return typ.Name()
}

// wartoscPozaZakresem sprawdza jedną wartość, napis albo tablicę napisów, wobec zakresu wyliczenia kontraktu; wartość pusta oraz wartość `null` przechodzą bez sprawdzenia, bo dla pola opcjonalnego znaczą brak wskazania.
func wartoscPozaZakresem(nazwa string, wartosc json.RawMessage, zakres []string) string {
	var napis string
	if err := json.Unmarshal(wartosc, &napis); err == nil {
		if napis == "" || naliscie(napis, zakres) {
			return ""
		}
		return "pole " + nazwa + ": wartość " + napis + " nie należy do wyliczenia kontraktu (" +
			strings.Join(zakres, ", ") + ")"
	}
	var lista []string
	if err := json.Unmarshal(wartosc, &lista); err != nil {
		return ""
	}
	for _, element := range lista {
		if element == "" || naliscie(element, zakres) {
			continue
		}
		return "pole " + nazwa + ": wartość " + element + " nie należy do wyliczenia kontraktu (" +
			strings.Join(zakres, ", ") + ")"
	}
	return ""
}

// naliscie odpowiada, czy podana wartość stoi w zakresie dopuszczalnych wartości wyliczenia kontraktu przekazanym jako parametr.
func naliscie(wartosc string, zakres []string) bool {
	for _, dopuszczalna := range zakres {
		if wartosc == dopuszczalna {
			return true
		}
	}
	return false
}

// zakresWyliczenia sprowadza wartości wyliczenia do napisów, bo kontrakt nadaje każdemu wyliczeniu własny typ napisowy, a treść żądania niesie sam napis.
func zakresWyliczenia[T ~string](wartosci []T) []string {
	zakres := make([]string, len(wartosci))
	for i, wartosc := range wartosci {
		zakres[i] = string(wartosc)
	}
	return zakres
}

// zakresyWyliczenKontraktu składa raz wykaz „wyliczenie kontraktu → dopuszczalne
// wartości" z sekcji wyliczeń kontraktu. Wykaz obejmuje sekcję w całości, a nie
// deklaracje narzędzi modelu, bo wykaz narzędzi z zamysłu pomija część komend
// i ich pola wyliczeniowe zostawałyby bez zakresu do sprawdzenia.
var zakresyWyliczenKontraktu = sync.OnceValue(func() map[string][]string {
	return map[string][]string{
		"EnvelopeStatus":                   zakresWyliczenia(shared.WartosciEnvelopeStatus()),
		"ChangeKind":                       zakresWyliczenia(shared.WartosciChangeKind()),
		"SessionStatus":                    zakresWyliczenia(shared.WartosciSessionStatus()),
		"WindowStatus":                     zakresWyliczenia(shared.WartosciWindowStatus()),
		"WindowRole":                       zakresWyliczenia(shared.WartosciWindowRole()),
		"ExecutionEnv":                     zakresWyliczenia(shared.WartosciExecutionEnv()),
		"PermissionMode":                   zakresWyliczenia(shared.WartosciPermissionMode()),
		"MessageRole":                      zakresWyliczenia(shared.WartosciMessageRole()),
		"MessageStatus":                    zakresWyliczenia(shared.WartosciMessageStatus()),
		"ChunkKind":                        zakresWyliczenia(shared.WartosciChunkKind()),
		"QueueStatus":                      zakresWyliczenia(shared.WartosciQueueStatus()),
		"QueueAction":                      zakresWyliczenia(shared.WartosciQueueAction()),
		"ProgressStatus":                   zakresWyliczenia(shared.WartosciProgressStatus()),
		"ConfigScope":                      zakresWyliczenia(shared.WartosciConfigScope()),
		"NavigationKind":                   zakresWyliczenia(shared.WartosciNavigationKind()),
		"ConfigAxis":                       zakresWyliczenia(shared.WartosciConfigAxis()),
		"SettingValueType":                 zakresWyliczenia(shared.WartosciSettingValueType()),
		"AccessPointKind":                  zakresWyliczenia(shared.WartosciAccessPointKind()),
		"AccessPointStatus":                zakresWyliczenia(shared.WartosciAccessPointStatus()),
		"AccessMode":                       zakresWyliczenia(shared.WartosciAccessMode()),
		"AccountKind":                      zakresWyliczenia(shared.WartosciAccountKind()),
		"IdentityLayer":                    zakresWyliczenia(shared.WartosciIdentityLayer()),
		"IdentityMode":                     zakresWyliczenia(shared.WartosciIdentityMode()),
		"LoopStopReason":                   zakresWyliczenia(shared.WartosciLoopStopReason()),
		"SessionConfigArea":                zakresWyliczenia(shared.WartosciSessionConfigArea()),
		"ReasoningEffort":                  zakresWyliczenia(shared.WartosciReasoningEffort()),
		"AccountSelection":                 zakresWyliczenia(shared.WartosciAccountSelection()),
		"ProviderTransport":                zakresWyliczenia(shared.WartosciProviderTransport()),
		"SystemPromptOrigin":               zakresWyliczenia(shared.WartosciSystemPromptOrigin()),
		"ToolPolicy":                       zakresWyliczenia(shared.WartosciToolPolicy()),
		"McpTransport":                     zakresWyliczenia(shared.WartosciMcpTransport()),
		"InteractionMode":                  zakresWyliczenia(shared.WartosciInteractionMode()),
		"InputFormat":                      zakresWyliczenia(shared.WartosciInputFormat()),
		"OutputFormat":                     zakresWyliczenia(shared.WartosciOutputFormat()),
		"CapabilitySupport":                zakresWyliczenia(shared.WartosciCapabilitySupport()),
		"AdapterSurface":                   zakresWyliczenia(shared.WartosciAdapterSurface()),
		"SessionConfigSource":              zakresWyliczenia(shared.WartosciSessionConfigSource()),
		"WorkingDirectoryDegradation":      zakresWyliczenia(shared.WartosciWorkingDirectoryDegradation()),
		"StudioDocumentFormat":             zakresWyliczenia(shared.WartosciStudioDocumentFormat()),
		"StudioOperationScope":             zakresWyliczenia(shared.WartosciStudioOperationScope()),
		"DiffHunkKind":                     zakresWyliczenia(shared.WartosciDiffHunkKind()),
		"ExportFormat":                     zakresWyliczenia(shared.WartosciExportFormat()),
		"WorkspaceProjectStatus":           zakresWyliczenia(shared.WartosciWorkspaceProjectStatus()),
		"MemoryEntryOrigin":                zakresWyliczenia(shared.WartosciMemoryEntryOrigin()),
		"AutomationStepKind":               zakresWyliczenia(shared.WartosciAutomationStepKind()),
		"AutomationTriggerKind":            zakresWyliczenia(shared.WartosciAutomationTriggerKind()),
		"AutomationExecutionStatus":        zakresWyliczenia(shared.WartosciAutomationExecutionStatus()),
		"AutomationDependencyKind":         zakresWyliczenia(shared.WartosciAutomationDependencyKind()),
		"ResearchSourceKind":               zakresWyliczenia(shared.WartosciResearchSourceKind()),
		"ResearchCredibility":              zakresWyliczenia(shared.WartosciResearchCredibility()),
		"ResearchFindingStatus":            zakresWyliczenia(shared.WartosciResearchFindingStatus()),
		"LibraryPreviewKind":               zakresWyliczenia(shared.WartosciLibraryPreviewKind()),
		"TranslationStatus":                zakresWyliczenia(shared.WartosciTranslationStatus()),
		"TranslationIssueKind":             zakresWyliczenia(shared.WartosciTranslationIssueKind()),
		"RoundtableFormat":                 zakresWyliczenia(shared.WartosciRoundtableFormat()),
		"RoundtableTurnStatus":             zakresWyliczenia(shared.WartosciRoundtableTurnStatus()),
		"ModeratorAction":                  zakresWyliczenia(shared.WartosciModeratorAction()),
		"DesignAssetKind":                  zakresWyliczenia(shared.WartosciDesignAssetKind()),
		"AssistantActionStatus":            zakresWyliczenia(shared.WartosciAssistantActionStatus()),
		"AssistantActionControl":           zakresWyliczenia(shared.WartosciAssistantActionControl()),
		"AssistantActivityKind":            zakresWyliczenia(shared.WartosciAssistantActivityKind()),
		"AssistantOrigin":                  zakresWyliczenia(shared.WartosciAssistantOrigin()),
		"TerminalShell":                    zakresWyliczenia(shared.WartosciTerminalShell()),
		"TerminalSessionStatus":            zakresWyliczenia(shared.WartosciTerminalSessionStatus()),
		"ProcessInitiator":                 zakresWyliczenia(shared.WartosciProcessInitiator()),
		"TerminalProcessStatus":            zakresWyliczenia(shared.WartosciTerminalProcessStatus()),
		"TreeNodeKind":                     zakresWyliczenia(shared.WartosciTreeNodeKind()),
		"GitActionKind":                    zakresWyliczenia(shared.WartosciGitActionKind()),
		"BuildStatus":                      zakresWyliczenia(shared.WartosciBuildStatus()),
		"ProblemSeverity":                  zakresWyliczenia(shared.WartosciProblemSeverity()),
		"LogLevel":                         zakresWyliczenia(shared.WartosciLogLevel()),
		"DiagnosticErrorStatus":            zakresWyliczenia(shared.WartosciDiagnosticErrorStatus()),
		"DiagnosticPriority":               zakresWyliczenia(shared.WartosciDiagnosticPriority()),
		"RecommendationStatus":             zakresWyliczenia(shared.WartosciRecommendationStatus()),
		"AppComponentKind":                 zakresWyliczenia(shared.WartosciAppComponentKind()),
		"AppArchitectureTemplate":          zakresWyliczenia(shared.WartosciAppArchitectureTemplate()),
		"AppWorkspaceLayer":                zakresWyliczenia(shared.WartosciAppWorkspaceLayer()),
		"AppDeployEnvironment":             zakresWyliczenia(shared.WartosciAppDeployEnvironment()),
		"AppDeployStrategy":                zakresWyliczenia(shared.WartosciAppDeployStrategy()),
		"AppDeployStatus":                  zakresWyliczenia(shared.WartosciAppDeployStatus()),
		"AppStageStatus":                   zakresWyliczenia(shared.WartosciAppStageStatus()),
		"AgentPermissionGroup":             zakresWyliczenia(shared.WartosciAgentPermissionGroup()),
		"MemoryLevel":                      zakresWyliczenia(shared.WartosciMemoryLevel()),
		"AgentVisibility":                  zakresWyliczenia(shared.WartosciAgentVisibility()),
		"AgentConnectorKind":               zakresWyliczenia(shared.WartosciAgentConnectorKind()),
		"IsolationContextKind":             zakresWyliczenia(shared.WartosciIsolationContextKind()),
		"IsolationTechnicalScope":          zakresWyliczenia(shared.WartosciIsolationTechnicalScope()),
		"IsolationLayer":                   zakresWyliczenia(shared.WartosciIsolationLayer()),
		"ComponentKind":                    zakresWyliczenia(shared.WartosciComponentKind()),
		"ExtensionKind":                    zakresWyliczenia(shared.WartosciExtensionKind()),
		"ExtensionOrigin":                  zakresWyliczenia(shared.WartosciExtensionOrigin()),
		"SubagentStatus":                   zakresWyliczenia(shared.WartosciSubagentStatus()),
		"TerminalOutputChannel":            zakresWyliczenia(shared.WartosciTerminalOutputChannel()),
		"MobileProcessControl":             zakresWyliczenia(shared.WartosciMobileProcessControl()),
		"AuthMethodKind":                   zakresWyliczenia(shared.WartosciAuthMethodKind()),
		"AuthChangeReason":                 zakresWyliczenia(shared.WartosciAuthChangeReason()),
		"AdvisorSelection":                 zakresWyliczenia(shared.WartosciAdvisorSelection()),
		"ImageTransformKind":               zakresWyliczenia(shared.WartosciImageTransformKind()),
		"ImageAdjustKind":                  zakresWyliczenia(shared.WartosciImageAdjustKind()),
		"MediaOperationKind":               zakresWyliczenia(shared.WartosciMediaOperationKind()),
		"MailProtocol":                     zakresWyliczenia(shared.WartosciMailProtocol()),
		"MailAccountSource":                zakresWyliczenia(shared.WartosciMailAccountSource()),
		"KnowledgeScope":                   zakresWyliczenia(shared.WartosciKnowledgeScope()),
		"ActorKind":                        zakresWyliczenia(shared.WartosciActorKind()),
		"SlashEntryKind":                   zakresWyliczenia(shared.WartosciSlashEntryKind()),
		"SessionToolSource":                zakresWyliczenia(shared.WartosciSessionToolSource()),
		"AgentAssignmentKind":              zakresWyliczenia(shared.WartosciAgentAssignmentKind()),
		"AgentPermissionAction":            zakresWyliczenia(shared.WartosciAgentPermissionAction()),
		"AgentModuleScope":                 zakresWyliczenia(shared.WartosciAgentModuleScope()),
		"AppArtifactKind":                  zakresWyliczenia(shared.WartosciAppArtifactKind()),
		"AppDataFlowKind":                  zakresWyliczenia(shared.WartosciAppDataFlowKind()),
		"AppEndpointMethod":                zakresWyliczenia(shared.WartosciAppEndpointMethod()),
		"AppEndpointStatus":                zakresWyliczenia(shared.WartosciAppEndpointStatus()),
		"AppExportFormat":                  zakresWyliczenia(shared.WartosciAppExportFormat()),
		"AppMilestoneStatus":               zakresWyliczenia(shared.WartosciAppMilestoneStatus()),
		"AppPackageFormat":                 zakresWyliczenia(shared.WartosciAppPackageFormat()),
		"AppPackageVisibility":             zakresWyliczenia(shared.WartosciAppPackageVisibility()),
		"AppPreviewStatus":                 zakresWyliczenia(shared.WartosciAppPreviewStatus()),
		"AppProductPlatform":               zakresWyliczenia(shared.WartosciAppProductPlatform()),
		"AppTimelineKind":                  zakresWyliczenia(shared.WartosciAppTimelineKind()),
		"AppValidationSeverity":            zakresWyliczenia(shared.WartosciAppValidationSeverity()),
		"ExtensionAuthKind":                zakresWyliczenia(shared.WartosciExtensionAuthKind()),
		"ExtensionBulkAction":              zakresWyliczenia(shared.WartosciExtensionBulkAction()),
		"ExtensionDefinitionFormat":        zakresWyliczenia(shared.WartosciExtensionDefinitionFormat()),
		"ExtensionHealthStatus":            zakresWyliczenia(shared.WartosciExtensionHealthStatus()),
		"ExtensionLifecycleAction":         zakresWyliczenia(shared.WartosciExtensionLifecycleAction()),
		"ExtensionPermissionScope":         zakresWyliczenia(shared.WartosciExtensionPermissionScope()),
		"ExtensionToolKind":                zakresWyliczenia(shared.WartosciExtensionToolKind()),
		"ExtensionTrustLevel":              zakresWyliczenia(shared.WartosciExtensionTrustLevel()),
		"ExtensionWebhookDirection":        zakresWyliczenia(shared.WartosciExtensionWebhookDirection()),
		"ProtocolFrameDirection":           zakresWyliczenia(shared.WartosciProtocolFrameDirection()),
		"ListenMode":                       zakresWyliczenia(shared.WartosciListenMode()),
		"ClipboardEntryKind":               zakresWyliczenia(shared.WartosciClipboardEntryKind()),
		"QueueItemStatus":                  zakresWyliczenia(shared.WartosciQueueItemStatus()),
		"QueueScope":                       zakresWyliczenia(shared.WartosciQueueScope()),
		"QueueBackoffKind":                 zakresWyliczenia(shared.WartosciQueueBackoffKind()),
		"AutomationLogLevel":               zakresWyliczenia(shared.WartosciAutomationLogLevel()),
		"AutomationAlertTrigger":           zakresWyliczenia(shared.WartosciAutomationAlertTrigger()),
		"AutomationTriggerCause":           zakresWyliczenia(shared.WartosciAutomationTriggerCause()),
		"OrchestrationGateRule":            zakresWyliczenia(shared.WartosciOrchestrationGateRule()),
		"BrowserMonitorStatus":             zakresWyliczenia(shared.WartosciBrowserMonitorStatus()),
		"BrowserNotifyChannel":             zakresWyliczenia(shared.WartosciBrowserNotifyChannel()),
		"BrowserFeedFormat":                zakresWyliczenia(shared.WartosciBrowserFeedFormat()),
		"BrowserConsoleLevel":              zakresWyliczenia(shared.WartosciBrowserConsoleLevel()),
		"BrowserScreenshotMode":            zakresWyliczenia(shared.WartosciBrowserScreenshotMode()),
		"BrowserImageFormat":               zakresWyliczenia(shared.WartosciBrowserImageFormat()),
		"BrowserDevicePreset":              zakresWyliczenia(shared.WartosciBrowserDevicePreset()),
		"BrowserDeviceOrientation":         zakresWyliczenia(shared.WartosciBrowserDeviceOrientation()),
		"BrowserWaitCondition":             zakresWyliczenia(shared.WartosciBrowserWaitCondition()),
		"BrowserNoteClassification":        zakresWyliczenia(shared.WartosciBrowserNoteClassification()),
		"BrowserMacroAction":               zakresWyliczenia(shared.WartosciBrowserMacroAction()),
		"BrowserDownloadStatus":            zakresWyliczenia(shared.WartosciBrowserDownloadStatus()),
		"BrowserDownloadAction":            zakresWyliczenia(shared.WartosciBrowserDownloadAction()),
		"BrowserArtifactKind":              zakresWyliczenia(shared.WartosciBrowserArtifactKind()),
		"BrowserTabState":                  zakresWyliczenia(shared.WartosciBrowserTabState()),
		"AssetContentDisposition":          zakresWyliczenia(shared.WartosciAssetContentDisposition()),
		"ImageVectorizeMode":               zakresWyliczenia(shared.WartosciImageVectorizeMode()),
		"ImageBlendMode":                   zakresWyliczenia(shared.WartosciImageBlendMode()),
		"DesignTokenTarget":                zakresWyliczenia(shared.WartosciDesignTokenTarget()),
		"DesignTokenKind":                  zakresWyliczenia(shared.WartosciDesignTokenKind()),
		"GitFileState":                     zakresWyliczenia(shared.WartosciGitFileState()),
		"GitDiffLineKind":                  zakresWyliczenia(shared.WartosciGitDiffLineKind()),
		"ConflictResolutionKind":           zakresWyliczenia(shared.WartosciConflictResolutionKind()),
		"SymbolNavigationKind":             zakresWyliczenia(shared.WartosciSymbolNavigationKind()),
		"RefactorKind":                     zakresWyliczenia(shared.WartosciRefactorKind()),
		"TestStatus":                       zakresWyliczenia(shared.WartosciTestStatus()),
		"DebugStatus":                      zakresWyliczenia(shared.WartosciDebugStatus()),
		"DebugStepKind":                    zakresWyliczenia(shared.WartosciDebugStepKind()),
		"BreakpointKind":                   zakresWyliczenia(shared.WartosciBreakpointKind()),
		"DataEngine":                       zakresWyliczenia(shared.WartosciDataEngine()),
		"SchemaNodeKind":                   zakresWyliczenia(shared.WartosciSchemaNodeKind()),
		"ContainerStatus":                  zakresWyliczenia(shared.WartosciContainerStatus()),
		"ContainerActionKind":              zakresWyliczenia(shared.WartosciContainerActionKind()),
		"ScanKind":                         zakresWyliczenia(shared.WartosciScanKind()),
		"ContextualOpKind":                 zakresWyliczenia(shared.WartosciContextualOpKind()),
		"ModelCallStatus":                  zakresWyliczenia(shared.WartosciModelCallStatus()),
		"ModelCallSpanKind":                zakresWyliczenia(shared.WartosciModelCallSpanKind()),
		"ModelCallQuality":                 zakresWyliczenia(shared.WartosciModelCallQuality()),
		"TelemetryFormat":                  zakresWyliczenia(shared.WartosciTelemetryFormat()),
		"UsageDimension":                   zakresWyliczenia(shared.WartosciUsageDimension()),
		"AlertRuleKind":                    zakresWyliczenia(shared.WartosciAlertRuleKind()),
		"AlertMetric":                      zakresWyliczenia(shared.WartosciAlertMetric()),
		"AlertComparison":                  zakresWyliczenia(shared.WartosciAlertComparison()),
		"AlertChannel":                     zakresWyliczenia(shared.WartosciAlertChannel()),
		"AlertTriggerStatus":               zakresWyliczenia(shared.WartosciAlertTriggerStatus()),
		"HealthProbeKind":                  zakresWyliczenia(shared.WartosciHealthProbeKind()),
		"HealthProbeStatus":                zakresWyliczenia(shared.WartosciHealthProbeStatus()),
		"LibraryFileStatus":                zakresWyliczenia(shared.WartosciLibraryFileStatus()),
		"LibraryRuleKind":                  zakresWyliczenia(shared.WartosciLibraryRuleKind()),
		"LibraryPreservationKind":          zakresWyliczenia(shared.WartosciLibraryPreservationKind()),
		"LibraryDuplicateKind":             zakresWyliczenia(shared.WartosciLibraryDuplicateKind()),
		"LibraryAuditAction":               zakresWyliczenia(shared.WartosciLibraryAuditAction()),
		"LibraryShareScope":                zakresWyliczenia(shared.WartosciLibraryShareScope()),
		"LibraryFieldKind":                 zakresWyliczenia(shared.WartosciLibraryFieldKind()),
		"LibraryThesaurusRelation":         zakresWyliczenia(shared.WartosciLibraryThesaurusRelation()),
		"LibraryRetentionAction":           zakresWyliczenia(shared.WartosciLibraryRetentionAction()),
		"LibraryPackageKind":               zakresWyliczenia(shared.WartosciLibraryPackageKind()),
		"LibrarySuggestionKind":            zakresWyliczenia(shared.WartosciLibrarySuggestionKind()),
		"LibraryDiffKind":                  zakresWyliczenia(shared.WartosciLibraryDiffKind()),
		"LibraryWebhookEvent":              zakresWyliczenia(shared.WartosciLibraryWebhookEvent()),
		"ResearchFindingKind":              zakresWyliczenia(shared.WartosciResearchFindingKind()),
		"ResearchFindingWeight":            zakresWyliczenia(shared.WartosciResearchFindingWeight()),
		"ResearchReadingState":             zakresWyliczenia(shared.WartosciResearchReadingState()),
		"ResearchAnnotationKind":           zakresWyliczenia(shared.WartosciResearchAnnotationKind()),
		"ResearchAnchorKind":               zakresWyliczenia(shared.WartosciResearchAnchorKind()),
		"ResearchDiscoveryMode":            zakresWyliczenia(shared.WartosciResearchDiscoveryMode()),
		"ResearchMonitorKind":              zakresWyliczenia(shared.WartosciResearchMonitorKind()),
		"ResearchAttachmentKind":           zakresWyliczenia(shared.WartosciResearchAttachmentKind()),
		"ResearchImportFormat":             zakresWyliczenia(shared.WartosciResearchImportFormat()),
		"ResearchCaptureMode":              zakresWyliczenia(shared.WartosciResearchCaptureMode()),
		"ResearchIdentifierKind":           zakresWyliczenia(shared.WartosciResearchIdentifierKind()),
		"ResearchSnowballDirection":        zakresWyliczenia(shared.WartosciResearchSnowballDirection()),
		"ResearchFactCheckVerdict":         zakresWyliczenia(shared.WartosciResearchFactCheckVerdict()),
		"ResearchCitationMode":             zakresWyliczenia(shared.WartosciResearchCitationMode()),
		"ResearchBibliographyScope":        zakresWyliczenia(shared.WartosciResearchBibliographyScope()),
		"ResearchBlockKind":                zakresWyliczenia(shared.WartosciResearchBlockKind()),
		"ResearchExportTarget":             zakresWyliczenia(shared.WartosciResearchExportTarget()),
		"RoundtableSpeechAct":              zakresWyliczenia(shared.WartosciRoundtableSpeechAct()),
		"RoundtableArgumentRelation":       zakresWyliczenia(shared.WartosciRoundtableArgumentRelation()),
		"RoundtableAnalysisKind":           zakresWyliczenia(shared.WartosciRoundtableAnalysisKind()),
		"RoundtableParticipantRole":        zakresWyliczenia(shared.WartosciRoundtableParticipantRole()),
		"RoundtableVoteMethod":             zakresWyliczenia(shared.WartosciRoundtableVoteMethod()),
		"RoundtableVoteStatus":             zakresWyliczenia(shared.WartosciRoundtableVoteStatus()),
		"RoundtableRatingKind":             zakresWyliczenia(shared.WartosciRoundtableRatingKind()),
		"RoundtableLeaderboardAlgorithm":   zakresWyliczenia(shared.WartosciRoundtableLeaderboardAlgorithm()),
		"RoundtableLeaderboardScope":       zakresWyliczenia(shared.WartosciRoundtableLeaderboardScope()),
		"RoundtableHandoffTarget":          zakresWyliczenia(shared.WartosciRoundtableHandoffTarget()),
		"RoundtableTranscriptFormat":       zakresWyliczenia(shared.WartosciRoundtableTranscriptFormat()),
		"RoundtableArgumentFormat":         zakresWyliczenia(shared.WartosciRoundtableArgumentFormat()),
		"StudioAuthor":                     zakresWyliczenia(shared.WartosciStudioAuthor()),
		"StudioChangeKind":                 zakresWyliczenia(shared.WartosciStudioChangeKind()),
		"StudioChangeDecision":             zakresWyliczenia(shared.WartosciStudioChangeDecision()),
		"StudioMergeSide":                  zakresWyliczenia(shared.WartosciStudioMergeSide()),
		"StudioPageOrientation":            zakresWyliczenia(shared.WartosciStudioPageOrientation()),
		"StudioPdfPageOperationKind":       zakresWyliczenia(shared.WartosciStudioPdfPageOperationKind()),
		"StudioIngestState":                zakresWyliczenia(shared.WartosciStudioIngestState()),
		"StudioOcrEngine":                  zakresWyliczenia(shared.WartosciStudioOcrEngine()),
		"StudioLayoutBlockKind":            zakresWyliczenia(shared.WartosciStudioLayoutBlockKind()),
		"StudioInputDeviceKind":            zakresWyliczenia(shared.WartosciStudioInputDeviceKind()),
		"TerminalScriptKind":               zakresWyliczenia(shared.WartosciTerminalScriptKind()),
		"TerminalTunnelKind":               zakresWyliczenia(shared.WartosciTerminalTunnelKind()),
		"TerminalTunnelStatus":             zakresWyliczenia(shared.WartosciTerminalTunnelStatus()),
		"TerminalKeyType":                  zakresWyliczenia(shared.WartosciTerminalKeyType()),
		"TerminalWatchStatus":              zakresWyliczenia(shared.WartosciTerminalWatchStatus()),
		"TerminalLintSeverity":             zakresWyliczenia(shared.WartosciTerminalLintSeverity()),
		"TranslationMemoryScope":           zakresWyliczenia(shared.WartosciTranslationMemoryScope()),
		"TranslationMemoryMaintenanceKind": zakresWyliczenia(shared.WartosciTranslationMemoryMaintenanceKind()),
		"ProofreadCheckKind":               zakresWyliczenia(shared.WartosciProofreadCheckKind()),
		"ProofreadSeverity":                zakresWyliczenia(shared.WartosciProofreadSeverity()),
		"ConsistencyFindingKind":           zakresWyliczenia(shared.WartosciConsistencyFindingKind()),
		"ApprovalStage":                    zakresWyliczenia(shared.WartosciApprovalStage()),
		"TranslationDocumentFormat":        zakresWyliczenia(shared.WartosciTranslationDocumentFormat()),
		"LayoutDifferenceKind":             zakresWyliczenia(shared.WartosciLayoutDifferenceKind()),
		"LocalizationResourceFormat":       zakresWyliczenia(shared.WartosciLocalizationResourceFormat()),
		"XliffVersion":                     zakresWyliczenia(shared.WartosciXliffVersion()),
		"SubtitleFormat":                   zakresWyliczenia(shared.WartosciSubtitleFormat()),
		"SubtitleTimingIssueKind":          zakresWyliczenia(shared.WartosciSubtitleTimingIssueKind()),
		"DubbingDurationSource":            zakresWyliczenia(shared.WartosciDubbingDurationSource()),
		"BatchOperationKind":               zakresWyliczenia(shared.WartosciBatchOperationKind()),
		"HandoffContent":                   zakresWyliczenia(shared.WartosciHandoffContent()),
		"HandoffStatus":                    zakresWyliczenia(shared.WartosciHandoffStatus()),
		"BridgeResultMode":                 zakresWyliczenia(shared.WartosciBridgeResultMode()),
		"TranslationArtifactKind":          zakresWyliczenia(shared.WartosciTranslationArtifactKind()),
		"TranslationStepKind":              zakresWyliczenia(shared.WartosciTranslationStepKind()),
		"GlossaryTermStatus":               zakresWyliczenia(shared.WartosciGlossaryTermStatus()),
		"TranslationExportVariant":         zakresWyliczenia(shared.WartosciTranslationExportVariant()),
		"WorkspaceTaskStatus":              zakresWyliczenia(shared.WartosciWorkspaceTaskStatus()),
		"WorkspaceTaskPriority":            zakresWyliczenia(shared.WartosciWorkspaceTaskPriority()),
		"WorkspaceAssigneeKind":            zakresWyliczenia(shared.WartosciWorkspaceAssigneeKind()),
		"WorkspaceDependencyKind":          zakresWyliczenia(shared.WartosciWorkspaceDependencyKind()),
		"WorkspaceCalendarSpan":            zakresWyliczenia(shared.WartosciWorkspaceCalendarSpan()),
		"WorkspaceEntityKind":              zakresWyliczenia(shared.WartosciWorkspaceEntityKind()),
		"WorkspaceGraphEdgeKind":           zakresWyliczenia(shared.WartosciWorkspaceGraphEdgeKind()),
		"WorkspaceExtractionMethod":        zakresWyliczenia(shared.WartosciWorkspaceExtractionMethod()),
		"WorkspaceActivityKind":            zakresWyliczenia(shared.WartosciWorkspaceActivityKind()),
		"DesignVectorNodeKind":             zakresWyliczenia(shared.WartosciDesignVectorNodeKind()),
		"DesignBooleanOp":                  zakresWyliczenia(shared.WartosciDesignBooleanOp()),
		"DesignShapeKind":                  zakresWyliczenia(shared.WartosciDesignShapeKind()),
		"DesignFillKind":                   zakresWyliczenia(shared.WartosciDesignFillKind()),
		"DesignGradientKind":               zakresWyliczenia(shared.WartosciDesignGradientKind()),
		"DesignStrokeCap":                  zakresWyliczenia(shared.WartosciDesignStrokeCap()),
		"DesignStrokeJoin":                 zakresWyliczenia(shared.WartosciDesignStrokeJoin()),
		"DesignVectorExportTarget":         zakresWyliczenia(shared.WartosciDesignVectorExportTarget()),
		"DesignLayoutDirection":            zakresWyliczenia(shared.WartosciDesignLayoutDirection()),
		"DesignLayoutAlign":                zakresWyliczenia(shared.WartosciDesignLayoutAlign()),
		"DesignConstraintAnchor":           zakresWyliczenia(shared.WartosciDesignConstraintAnchor()),
		"DesignPrototypeTrigger":           zakresWyliczenia(shared.WartosciDesignPrototypeTrigger()),
		"DesignPrototypeTransition":        zakresWyliczenia(shared.WartosciDesignPrototypeTransition()),
		"DesignColorHarmony":               zakresWyliczenia(shared.WartosciDesignColorHarmony()),
		"DesignColorVision":                zakresWyliczenia(shared.WartosciDesignColorVision()),
		"DesignColorSpace":                 zakresWyliczenia(shared.WartosciDesignColorSpace()),
		"DesignIconSpriteKind":             zakresWyliczenia(shared.WartosciDesignIconSpriteKind()),
		"DesignPrintColorSpace":            zakresWyliczenia(shared.WartosciDesignPrintColorSpace()),
		"DesignPrintStandard":              zakresWyliczenia(shared.WartosciDesignPrintStandard()),
		"DesignPreflightSeverity":          zakresWyliczenia(shared.WartosciDesignPreflightSeverity()),
		"DesignChartKind":                  zakresWyliczenia(shared.WartosciDesignChartKind()),
		"DesignDiagramKind":                zakresWyliczenia(shared.WartosciDesignDiagramKind()),
		"DesignTemplateKind":               zakresWyliczenia(shared.WartosciDesignTemplateKind()),
		"DesignPhotoResampleFilter":        zakresWyliczenia(shared.WartosciDesignPhotoResampleFilter()),
		"DesignPhotoFilter":                zakresWyliczenia(shared.WartosciDesignPhotoFilter()),
		"DesignPhotoBlendMode":             zakresWyliczenia(shared.WartosciDesignPhotoBlendMode()),
		"DesignPhotoRetouchMode":           zakresWyliczenia(shared.WartosciDesignPhotoRetouchMode()),
		"DesignPhotoComputeRoute":          zakresWyliczenia(shared.WartosciDesignPhotoComputeRoute()),
		"DesignPrintBinding":               zakresWyliczenia(shared.WartosciDesignPrintBinding()),
		"StudioTextAlign":                  zakresWyliczenia(shared.WartosciStudioTextAlign()),
		"StudioUnderlineStyle":             zakresWyliczenia(shared.WartosciStudioUnderlineStyle()),
		"StudioTextEffect":                 zakresWyliczenia(shared.WartosciStudioTextEffect()),
		"StudioLineSpacingRule":            zakresWyliczenia(shared.WartosciStudioLineSpacingRule()),
		"StudioTabKind":                    zakresWyliczenia(shared.WartosciStudioTabKind()),
		"StudioTabLeader":                  zakresWyliczenia(shared.WartosciStudioTabLeader()),
		"StudioBorderStyle":                zakresWyliczenia(shared.WartosciStudioBorderStyle()),
		"StudioCaseTransform":              zakresWyliczenia(shared.WartosciStudioCaseTransform()),
		"StudioStyleKind":                  zakresWyliczenia(shared.WartosciStudioStyleKind()),
		"StudioPaperKind":                  zakresWyliczenia(shared.WartosciStudioPaperKind()),
		"StudioSectionStart":               zakresWyliczenia(shared.WartosciStudioSectionStart()),
		"StudioBreakKind":                  zakresWyliczenia(shared.WartosciStudioBreakKind()),
		"StudioHeaderScope":                zakresWyliczenia(shared.WartosciStudioHeaderScope()),
		"StudioPageNumberFormat":           zakresWyliczenia(shared.WartosciStudioPageNumberFormat()),
		"StudioWatermarkKind":              zakresWyliczenia(shared.WartosciStudioWatermarkKind()),
		"StudioListKind":                   zakresWyliczenia(shared.WartosciStudioListKind()),
		"StudioListNumberFormat":           zakresWyliczenia(shared.WartosciStudioListNumberFormat()),
		"StudioBulletSource":               zakresWyliczenia(shared.WartosciStudioBulletSource()),
		"StudioTableStructureOp":           zakresWyliczenia(shared.WartosciStudioTableStructureOp()),
		"StudioTableConvert":               zakresWyliczenia(shared.WartosciStudioTableConvert()),
		"StudioVerticalAlign":              zakresWyliczenia(shared.WartosciStudioVerticalAlign()),
		"StudioObjectKind":                 zakresWyliczenia(shared.WartosciStudioObjectKind()),
		"StudioObjectSource":               zakresWyliczenia(shared.WartosciStudioObjectSource()),
		"StudioTextWrap":                   zakresWyliczenia(shared.WartosciStudioTextWrap()),
		"StudioAnchorKind":                 zakresWyliczenia(shared.WartosciStudioAnchorKind()),
		"StudioShapeKind":                  zakresWyliczenia(shared.WartosciStudioShapeKind()),
		"StudioApparatusKind":              zakresWyliczenia(shared.WartosciStudioApparatusKind()),
		"StudioFieldKind":                  zakresWyliczenia(shared.WartosciStudioFieldKind()),
		"StudioTemplateFieldKind":          zakresWyliczenia(shared.WartosciStudioTemplateFieldKind()),
		"StudioExportFormat":               zakresWyliczenia(shared.WartosciStudioExportFormat()),
		"StudioImportFormat":               zakresWyliczenia(shared.WartosciStudioImportFormat()),
		"StudioLockScope":                  zakresWyliczenia(shared.WartosciStudioLockScope()),
		"StudioActionKind":                 zakresWyliczenia(shared.WartosciStudioActionKind()),
		"StudioActionState":                zakresWyliczenia(shared.WartosciStudioActionState()),
		"StudioMarkupKind":                 zakresWyliczenia(shared.WartosciStudioMarkupKind()),
		"StudioMarkupState":                zakresWyliczenia(shared.WartosciStudioMarkupState()),
		"StudioSurfaceMode":                zakresWyliczenia(shared.WartosciStudioSurfaceMode()),
		"StudioSplitOrientation":           zakresWyliczenia(shared.WartosciStudioSplitOrientation()),
		"StudioZoomPreset":                 zakresWyliczenia(shared.WartosciStudioZoomPreset()),
		"StudioRulerUnit":                  zakresWyliczenia(shared.WartosciStudioRulerUnit()),
		"StudioScrollMode":                 zakresWyliczenia(shared.WartosciStudioScrollMode()),
		"StudioViewMode":                   zakresWyliczenia(shared.WartosciStudioViewMode()),
		"StudioPasteMode":                  zakresWyliczenia(shared.WartosciStudioPasteMode()),
		"StudioBackupReason":               zakresWyliczenia(shared.WartosciStudioBackupReason()),
		"StudioVersionSeries":              zakresWyliczenia(shared.WartosciStudioVersionSeries()),
		"StudioProvenanceKind":             zakresWyliczenia(shared.WartosciStudioProvenanceKind()),
		"StudioAgentConflictPolicy":        zakresWyliczenia(shared.WartosciStudioAgentConflictPolicy()),
		"StudioTaskKind":                   zakresWyliczenia(shared.WartosciStudioTaskKind()),
		"StudioTaskState":                  zakresWyliczenia(shared.WartosciStudioTaskState()),
		"StudioPlanState":                  zakresWyliczenia(shared.WartosciStudioPlanState()),
		"StudioAgentSlotState":             zakresWyliczenia(shared.WartosciStudioAgentSlotState()),
		"AodEventClass":                    zakresWyliczenia(shared.WartosciAodEventClass()),
		"AodMuteKind":                      zakresWyliczenia(shared.WartosciAodMuteKind()),
		"AodMuteScope":                     zakresWyliczenia(shared.WartosciAodMuteScope()),
		"NotificationClass":                zakresWyliczenia(shared.WartosciNotificationClass()),
		"NotificationWeight":               zakresWyliczenia(shared.WartosciNotificationWeight()),
		"NotificationState":                zakresWyliczenia(shared.WartosciNotificationState()),
		"NotificationSourceKind":           zakresWyliczenia(shared.WartosciNotificationSourceKind()),
		"NotificationDelivery":             zakresWyliczenia(shared.WartosciNotificationDelivery()),
		"BrowserAccessibilityLevel":        zakresWyliczenia(shared.WartosciBrowserAccessibilityLevel()),
		"BrowserAccessibilityStandard":     zakresWyliczenia(shared.WartosciBrowserAccessibilityStandard()),
		"AppPerformanceFormFactor":         zakresWyliczenia(shared.WartosciAppPerformanceFormFactor()),
		"SessionExportFormat":              zakresWyliczenia(shared.WartosciSessionExportFormat()),
	}
})

// bladZgodnosciZKontraktem nazywa niezgodność żądania z kontraktem. Komenda
// stoi w treści odmowy, bo odmowa idzie do Operatora oderwana od żądania.
func bladZgodnosciZKontraktem(komenda shared.MessageType, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		string(komenda)+": "+powod))
}

// bladWyliczeniaBezZakresu nazywa rozjazd wykazu bramy z kontraktem: pole niesie
// wyliczenie, którego wykaz nie zna, więc zakresu nie ma czym sprawdzić. Odmowa
// zamiast przepuszczenia, bo przepuszczone pole rozstrzyga dopiero CHECK
// sterownika bazy, a rozjazd bez odmowy nie zostawia po sobie śladu.
func bladWyliczeniaBezZakresu(komenda shared.MessageType, pole, wyliczenie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		string(komenda)+": pole "+pole+": wyliczenie "+wyliczenie+
			" bez zakresu w wykazie bramy kontraktu"))
}
