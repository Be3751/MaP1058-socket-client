package adapter

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Be3751/MaP1058-socket-client/internal/model"
	"github.com/Be3751/MaP1058-socket-client/internal/parser"
	"github.com/Be3751/MaP1058-socket-client/internal/socket"
	"io"
)

type BinAdapter interface {
	// WriteRawSignal AD値を受信する
	WriteRawSignal(ctx context.Context, stg *model.Setting) error
}

func NewBinAdapter(c socket.Conn, p parser.Parser, w io.ReadWriteSeeker) BinAdapter {
	return &binAdapter{
		Conn:   c,
		Parser: p,
		File:   w,
	}
}

type binAdapter struct {
	Conn   socket.Conn
	Parser parser.Parser
	File   io.ReadWriteSeeker
}

const (
	bufferSize = 10
)

func (a *binAdapter) WriteRawSignal(ctx context.Context, stg *model.Setting) (err error) {
	csvWriter := csv.NewWriter(a.File)
	defer func() {
		csvWriter.Flush()
		if hereErr := csvWriter.Error(); hereErr != nil {
			err = fmt.Errorf("%s: failed to flush csv writer: %w", err.Error(), hereErr)
		}
	}()
	if err = csvWriter.Write(model.SignalHeader()); err != nil {
		return fmt.Errorf("failed to write header to csv: %w", err)
	}

	buf := make([][]string, bufferSize)
	var timeReceived int
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			signals, err := a.receiveAD()
			if err != nil {
				return fmt.Errorf("failed to receive AD values: %w", err)
			}
			if err = signals.SetMeasurements(stg.Calibration, stg.AnalysisType); err != nil {
				return fmt.Errorf("failed to set measurements: %w", err)
			}
			buf = append(buf, signals.ToRecords(timeReceived)...)
			if err = a.sendACK(); err != nil {
				return err
			}
			timeReceived++
			if timeReceived == bufferSize {
				if err = a.writeRecords(csvWriter, buf); err != nil {
					return fmt.Errorf("failed to write raw signal records to csv: %w", err)
				}
				buf = make([][]string, bufferSize)
				timeReceived = 0
			}
		}
	}
}

func (a *binAdapter) receiveAD() (*model.Signals, error) {
	rawBytes := make([]byte, model.NumTotalBytes)
	n, err := a.Conn.Read(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to receive binary data: %w", err)
	}
	if n == 0 {
		return nil, errors.New("received 0 byte")
	}

	s, err := a.Parser.ToSignals(rawBytes[:n])
	if _, ok := err.(*parser.FailureSumCheckError); ok {
		if err := a.sendNAK(); err != nil {
			return nil, fmt.Errorf("%s, and failed to send NAK to the server", err.Error())
		}
	} else {
		return nil, fmt.Errorf("failed to parse binary data to Signals: %w", err)
	}
	return s, nil
}

func (a *binAdapter) sendACK() error {
	_, err := a.Conn.Write([]byte{0x06})
	if err != nil {
		return fmt.Errorf("failed to write connection ACK: %w", err)
	}
	return nil
}

func (a *binAdapter) sendNAK() error {
	_, err := a.Conn.Write([]byte{0x15})
	if err != nil {
		return fmt.Errorf("failed to write connection NAK: %w", err)
	}
	return nil
}

func (a *binAdapter) writeRecords(w *csv.Writer, records [][]string) error {
	err := w.WriteAll(records)
	if err != nil {
		return fmt.Errorf("failed to write records to csv: %w", err)
	}
	return nil
}
