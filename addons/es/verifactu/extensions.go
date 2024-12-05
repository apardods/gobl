package verifactu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for Verifactu
const (
	ExtKeyDocType        cbc.Key = "es-verifactu-doc-type"
	ExtKeyOpClass        cbc.Key = "es-verifactu-op-class"
	ExtKeyCorrectionType cbc.Key = "es-verifactu-correction-type"
	ExtKeyExempt         cbc.Key = "es-verifactu-exempt"
	ExtKeyIVAIGICRegime  cbc.Key = "es-verifactu-iva-igic-regime"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Verifactu Invoice Type Code - L2",
			i18n.ES: "Código de Tipo de Factura de Verifactu - L2",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Invoice type code used to identify the type of invoice being sent.
				Source: VeriFactu Ministerial Order:
				 * https://www.boe.es/diario_boe/txt.php?id=BOE-A-2024-22138
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "F1",
				Name: i18n.String{
					i18n.EN: "Invoice (Article 6, 7.2 and 7.3 of RD 1619/2012)",
					i18n.ES: "Factura (Art. 6, 7.2 y 7.3 del RD 1619/2012)",
				},
			},
			{
				Value: "F2",
				Name: i18n.String{
					i18n.EN: "Simplified invoice (Article 6.1.d) of RD 1619/2012)",
					i18n.ES: "Factura Simplificada (Art. 6.1.d) del RD 1619/2012)",
				},
			},
			{
				Value: "F3",
				Name: i18n.String{
					i18n.EN: "Invoice issued as a replacement for simplified invoices that have been billed and declared.",
					i18n.ES: "Factura emitida en sustitución de facturas simplificadas facturadas y declaradas.",
				},
			},
			{
				Value: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: error based on law and Article 80 One, Two and Six LIVA",
					i18n.ES: "Factura rectificativa: error fundado en derecho y Art. 80 Uno, Dos y Seis LIVA",
				},
			},
			{
				Value: "R2",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80.3",
					i18n.EN: "Rectified invoice: error based on law and Article 80.3",
				},
			},
			{
				Value: "R3",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80.4",
					i18n.EN: "Rectified invoice: error based on law and Article 80.4",
				},
			},
			{
				Value: "R4",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: Resto",
					i18n.EN: "Rectified invoice: Other",
				},
			},
			{
				Value: "R5",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: facturas simplificadas",
					i18n.EN: "Rectified invoice: simplified invoices",
				},
			},
		},
	},
	{
		Key: ExtKeyCorrectionType,
		Name: i18n.String{
			i18n.EN: "Verifactu Correction Type Code - L3",
			i18n.ES: "Código de Tipo de Corrección de Verifactu - L3",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Correction type code used to identify the type of correction being made.
				This value will be determined automatically according to the invoice type.
				Corrective invoices will be marked as "S", while credit and debit notes as "I".
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "S",
				Name: i18n.String{
					i18n.EN: "Substitution",
					i18n.ES: "Por Sustitución",
				},
			},
			{
				Value: "I",
				Name: i18n.String{
					i18n.EN: "Differences",
					i18n.ES: "Por Diferencias",
				},
			},
		},
	},
	{
		Key: ExtKeyOpClass,
		Name: i18n.String{
			i18n.EN: "Verifactu Operation Classification/Exemption Code - L9",
			i18n.ES: "Código de Clasificación/Exención de Impuesto de Verifactu - L9",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Operation classification code used to identify if taxes should be applied to the line.
				Source: VeriFactu Ministerial Order:
				 * https://www.boe.es/diario_boe/txt.php?id=BOE-A-2024-22138
				For details on how best to use and apply these and other codes, see the
				AEAT FAQ:
				 * https://sede.agenciatributaria.gob.es/Sede/impuestos-tasas/iva/iva-libros-registro-iva-traves-aeat/preguntas-frecuentes/3-libro-registro-facturas-expedidas.html?faqId=b5556c3d02bc9510VgnVCM100000dc381e0aRCRD
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "S1",
				Name: i18n.String{
					i18n.EN: "Subject and Not Exempt - Without reverse charge",
					i18n.ES: "Operación Sujeta y No exenta - Sin inversión del sujeto pasivo",
				},
			},
			{
				Value: "S2",
				Name: i18n.String{
					i18n.EN: "Subject and Not Exempt - With reverse charge",
					i18n.ES: "Operación Sujeta y No exenta - Con Inversión del sujeto pasivo",
				},
			},
			{
				Value: "N1",
				Name: i18n.String{
					i18n.EN: "Not Subject - Articles 7, 14, others",
					i18n.ES: "Operación No Sujeta artículo 7, 14, otros",
				},
			},
			{
				Value: "N2",
				Name: i18n.String{
					i18n.EN: "Not Subject - Due to location rules",
					i18n.ES: "Operación No Sujeta por Reglas de localización",
				},
			},
		},
	},
	{
		Key: ExtKeyExempt,
		Name: i18n.String{
			i18n.EN: "Verifactu Exemption Code - L10",
			i18n.ES: "Código de Exención de Impuesto de Verifactu - L10",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Exemption code used to explain why the operation is exempt from taxes.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 20. Exemptions in internal operations.",
					i18n.ES: "Exenta: por el artículo 20. Exenciones en operaciones interiores.",
				},
			},
			{
				Value: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 21. Exemptions in exports of goods.",
					i18n.ES: "Exenta: por el artículo 21. Exenciones en las exportaciones de bienes.",
				},
			},
			{
				Value: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 22. Exemptions in operations asimilated to exports.",
					i18n.ES: "Exenta: por el artículo 22. Exenciones en las operaciones asimiladas a las exportaciones.",
				},
			},
			{
				Value: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Articles 23 and 24. Exemptions related to temporary deposit, customs and fiscal regimes, and other situations.",
					i18n.ES: "Exenta: por el artículos 23 y 24. Exenciones relativas a las situaciones de depósito temporal, regímenes aduaneros y fiscales, y otras situaciones.",
				},
			},
			{
				Value: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 25. Exemptions in the delivery of goods destined to another Member State.",
					i18n.ES: "Exenta: por el artículo 25. Exenciones en las entregas de bienes destinados a otro Estado miembro.",
				},
			},
			{
				Value: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to other reasons",
					i18n.ES: "Exenta: por otra causa",
				},
			},
		},
	},
	{
		Key: ExtKeyIVAIGICRegime,
		Name: i18n.String{
			i18n.EN: "IVA/IGIC Regime Code - L8A/B",
			i18n.ES: "Código de Régimen de IVA/IGIC - L8A/B",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code list used to identify the type of VAT or IGIC regime applied to the operation. This list combines lists 8A, which include values for IVA and 8B, which include values for IGIC.
			`),
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "01",
				Name: i18n.String{
					i18n.EN: "General regime operation",
					i18n.ES: "Operación de régimen general",
				},
			},
			{
				Value: "02",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.ES: "Exportación",
				},
			},
			{
				Value: "03",
				Name: i18n.String{
					i18n.EN: "Special regime for used goods, art objects, antiques and collectibles",
					i18n.ES: "Operaciones a las que se aplique el régimen especial de bienes usados, objetos de arte, antigüedades y objetos de colección",
				},
			},
			{
				Value: "04",
				Name: i18n.String{
					i18n.EN: "Special regime for investment gold",
					i18n.ES: "Régimen especial del oro de inversión",
				},
			},
			{
				Value: "05",
				Name: i18n.String{
					i18n.EN: "Special regime for travel agencies",
					i18n.ES: "Régimen especial de las agencias de viajes",
				},
			},
			{
				Value: "06",
				Name: i18n.String{
					i18n.EN: "Special regime for IVA/IGIC groups (Advanced Level)",
					i18n.ES: "Régimen especial grupo de entidades en IVA/IGIC (Nivel Avanzado)",
				},
			},
			{
				Value: "07",
				Name: i18n.String{
					i18n.EN: "Special cash accounting regime",
					i18n.ES: "Régimen especial del criterio de caja",
				},
			},
			{
				Value: "08",
				Name: i18n.String{
					i18n.EN: "Operations subject to a different regime",
					i18n.ES: "Operaciones sujetas a un régimen diferente",
				},
			},
			{
				Value: "09",
				Name: i18n.String{
					i18n.EN: "Billing of travel agency services acting as mediators in name and on behalf of others",
					i18n.ES: "Facturación de las prestaciones de servicios de agencias de viaje que actúan como mediadoras en nombre y por cuenta ajena",
				},
			},
			{
				Value: "10",
				Name: i18n.String{
					i18n.EN: "Collection of professional fees or rights on behalf of third parties",
					i18n.ES: "Cobros por cuenta de terceros de honorarios profesionales o de derechos derivados de la propiedad industrial",
				},
			},
			{
				Value: "11",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations",
					i18n.ES: "Operaciones de arrendamiento de local de negocio",
				},
			},
			{
				Value: "14",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT accrual in work certifications for Public Administration",
					i18n.ES: "Factura con IVA pendiente de devengo en certificaciones de obra cuyo destinatario sea una Administración Pública",
				},
			},
			{
				Value: "15",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT accrual in successive tract operations",
					i18n.ES: "Factura con IVA pendiente de devengo en operaciones de tracto sucesivo",
				},
			},
			{
				Value: "17",
				Name: i18n.String{
					i18n.EN: "Operation under OSS and IOSS regimes",
					i18n.ES: "Operación acogida a alguno de los regímenes previstos en el capítulo XI del título IX (OSS e IOSS)",
				},
			},
			{
				Value: "18",
				Name: i18n.String{
					i18n.EN: "Equivalence surcharge (IVA) / Special regime for small traders or retailers (IGIC)",
					i18n.ES: "Recargo de equivalencia (IVA) / Régimen especial del pequeño comerciante o minorista (IGIC)",
				},
			},
			{
				Value: "19",
				Name: i18n.String{
					i18n.EN: "Operations included in the Special Regime for Agriculture, Livestock and Fisheries",
					i18n.ES: "Operaciones de actividades incluidas en el Régimen Especial de Agricultura, Ganadería y Pesca (REAGYP)",
				},
			},
			{
				Value: "20",
				Name: i18n.String{
					i18n.EN: "Simplified regime (IVA)",
					i18n.ES: "Régimen simplificado (IVA)",
				},
			},
		},
	},
}
