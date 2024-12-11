package models

import "time"

// CGDArchivos representa la estructura de la tabla CGD_ARCHIVO.
type CGDArchivos struct {
	IDArchivo                   int64     `json:"id_archivo" gorm:"type:numeric(16);primaryKey"`
	NombreArchivo               string    `json:"nombre_archivo" gorm:"type:varchar(100);not null"`
	PlataformaOrigen            string    `json:"plataforma_origen" gorm:"type:char(2);not null"`
	TipoArchivo                 string    `json:"tipo_archivo" gorm:"type:char(2);not null"`
	ConsecutivoPlataformaOrigen int16     `json:"consecutivo_plataforma_origen" gorm:"type:smallint;not null"`
	FechaNombreArchivo          string    `json:"fecha_nombre_archivo" gorm:"type:char(8);not null"`
	FechaRegistroResumen        string    `json:"fecha_registro_resumen" gorm:"type:char(14)"`
	NroTotalRegistros           int64     `json:"nro_total_registros" gorm:"type:numeric(9)"`
	NroRegistrosError           int64     `json:"nro_registros_error" gorm:"type:numeric(9)"`
	NroRegistrosValidos         int64     `json:"nro_registros_validos" gorm:"type:numeric(9)"`
	Estado                      string    `json:"estado" gorm:"type:varchar(50);not null"`
	FechaRecepcion              time.Time `json:"fecha_recepcion" gorm:"type:timestamp;not null"`
	FechaCiclo                  time.Time `json:"fecha_ciclo" gorm:"type:date;not null"`
	ContadorIntentosCargue      int16     `json:"contador_intentos_cargue" gorm:"type:smallint;not null"`
	ContadorIntentosGeneracion  int16     `json:"contador_intentos_generacion" gorm:"type:smallint;not null"`
	ContadorIntentosEmpaquetado int16     `json:"contador_intentos_empaquetado" gorm:"type:smallint;not null"`
	ACGFechaGeneracion          time.Time `json:"acg_fecha_generacion" gorm:"type:timestamp"`
	ACGConsecutivo              int64     `json:"acg_consecutivo" gorm:"type:numeric(4)"`
	ACGNombreArchivo            string    `json:"acg_nombre_archivo" gorm:"type:varchar(100)"`
	ACGRegistroEncabezado       string    `json:"acg_registro_encabezado" gorm:"type:varchar(200)"`
	ACGRegistroResumen          string    `json:"acg_registro_resumen" gorm:"type:varchar(200)"`
	ACGTotalTx                  int64     `json:"acg_total_tx" gorm:"type:numeric(9)"`
	ACGMontoTotalTx             float64   `json:"acg_monto_total_tx" gorm:"type:decimal(19,2)"`
	ACGTotalTxDebito            int64     `json:"acg_total_tx_debito" gorm:"type:numeric(9)"`
	ACGMontoTotalTxDebito       float64   `json:"acg_monto_total_tx_debito" gorm:"type:decimal(19,2)"`
	ACGTotalTxReverso           int64     `json:"acg_total_tx_reverso" gorm:"type:numeric(9)"`
	ACGMontoTotalTxReverso      float64   `json:"acg_monto_total_tx_reverso" gorm:"type:decimal(19,2)"`
	ACGTotalTxReintegro         int64     `json:"acg_total_tx_reintegro" gorm:"type:numeric(9)"`
	ACGMontoTotalTxReintegro    float64   `json:"acg_monto_total_tx_reintegro" gorm:"type:decimal(19,2)"`
	AnulacionNombreArchivo      string    `json:"anulacion_nombre_archivo" gorm:"type:varchar(100)"`
	AnulacionJustificacion      string    `json:"anulacion_justificacion" gorm:"type:varchar(4000)"`
	AnulacionFechaAnulacion     time.Time `json:"anulacion_fecha_anulacion" gorm:"type:timestamp"`
	GAWRtaTransEstado           string    `json:"gaw_rta_trans_estado" gorm:"type:varchar(50)"`
	GAWRtaTransCodigo           string    `json:"gaw_rta_trans_codigo" gorm:"type:varchar(4)"`
	GAWRtaTransDetalle          string    `json:"gaw_rta_trans_detalle" gorm:"type:varchar(1000)"`
	IDConsolidado               int64     `json:"id_consolidado" gorm:"type:numeric(14);foreignkey:IDConsolidado"`
	CodigoError                 string    `json:"codigo_error" gorm:"type:varchar(30);foreignkey:CodigoError"`
	DetalleError                string    `json:"detalle_error" gorm:"type:varchar(2000)"`
}

func (CGDArchivos) TableName() string {
	return "cgd_archivos"
}

// CGDArchivoEstados representa la estructura de la tabla CGD_ARCHIVO_ESTADO.
type CGDArchivoEstados struct {
	IDArchivo         int64     `json:"id_archivo" gorm:"type:numeric(16);primaryKey;foreignKey:IDArchivo"`
	EstadoInicial     string    `json:"estado_inicial" gorm:"type:varchar(50)"`
	EstadoFinal       string    `json:"estado_final" gorm:"type:varchar(50);primaryKey"`
	FechaCambioEstado time.Time `json:"fecha_cambio_estado" gorm:"type:timestamp;primaryKey;not null;autoCreateTime(6)"`
}

func (CGDArchivoEstados) TableName() string {
	return "cgd_archivo_estados"
}
